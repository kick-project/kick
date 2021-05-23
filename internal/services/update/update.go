package update

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/kick-project/kick/internal/resources/client"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/logger"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/model"
	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Update build metadata
type Update struct {
	client      *client.Client
	configFile  *config.File       `validate:"required"`
	orm         *gorm.DB           `validate:"required"`
	log         logger.OutputIface `validate:"required"`
	metadataDir string             `validate:"required"`
}

// Options constructor options
type Options struct {
	Client      *client.Client     `validate:"required"`
	ConfigFile  *config.File       `validate:"required"`
	ORM         *gorm.DB           `validate:"required"`
	Log         logger.OutputIface `validate:"required"`
	MetadataDir string             `validate:"required"`
}

// New constructor
func New(opts *Options) *Update {
	return &Update{
		client:      opts.Client,
		configFile:  opts.ConfigFile,
		orm:         opts.ORM,
		log:         opts.Log,
		metadataDir: opts.MetadataDir,
	}
}

// Build metadata. Conf defaults to globals.Config if Conf is nil.
func (m *Update) Build() error {
	conf := m.configFile

	c := workers{
		client: m.client,
		wait:   &sync.WaitGroup{},
		log:    m.log,
	}

	churl := make(chan string, 64)
	chtemplates := make(chan *Template, 64)
	c.concurClones(6, churl, chtemplates)
	c.concurInserts(m.orm, chtemplates)

	for _, url := range conf.RepoURLs {
		c.wait.Add(1)
		churl <- url
	}

	// Wait for all all processing to finish
	c.wait.Wait()

	return nil
}

type workers struct {
	client *client.Client
	log    logger.OutputIface
	wait   *sync.WaitGroup
}

// concurClones concurrent cloning of git repositories.
// where num is the number of concurrent downloads, churl is a string url and tchan is a channel of resulting templates.
func (c *workers) concurClones(num int, churl <-chan string, tchan chan<- *Template) {
	for i := 0; i < num; i++ {
		go func() {
			for {
				url, ok := <-churl
				switch {
				case !ok:
					return
				default:
					c.processURL(url, tchan)
					c.wait.Done()
				}
			}
		}()
	}
}

func (c *workers) processURL(url string, chtemplate chan<- *Template) {
	p, err := c.client.GetTemplate(url, "")
	if errs.LogF("error: cloning repository: %w: skipping %s", err, url) {
		return
	}
	localpath := p.Path()

	mpath := filepath.Clean(fmt.Sprintf("%s/repo.yml", localpath))
	if errs.LogF("error: can not open %s: %w: skipping %s", mpath, err, url) {
		return
	}

	repo := &Repo{URL: url}
	err = repo.Load(mpath)
	if errs.LogF("error: %w: skipping %s\n", err, url) {
		return
	}

	paths, err := filepath.Glob(filepath.Clean(fmt.Sprintf("%s/templates/*.yml", localpath)))
	if errs.LogF("error: getting a lists of paths: %w: skipping %s\n", err, url) {
		return
	}

	for _, curpath := range paths {
		t := &Template{}
		err := t.Load(curpath)
		if errs.LogF("error: loading template metadata from %s: %w: skipping", curpath, err) {
			continue
		}
		t.Repo = *repo
		c.wait.Add(1)
		chtemplate <- t
	}
}

// concurInserts populates the database
// where num is the number of concurrent routines
// and ch is the channel to read templates from.
func (c *workers) concurInserts(orm *gorm.DB, ch <-chan *Template) {
	go func() {
		for {
			t, ok := <-ch
			switch {
			case !ok:
				return
			default:
				c.insert(orm, t)
				c.wait.Done()
			}
		}
	}()
}

func (c *workers) insert(orm *gorm.DB, t *Template) {
	modRepo := model.Repo{
		Name: t.Repo.Name,
		URL:  t.Repo.URL,
		Desc: t.Repo.Description,
	}
	result := orm.Clauses(clause.Insert{Modifier: "OR IGNORE"}).Create(&modRepo)
	if result.RowsAffected != 1 {
		result = orm.First(&modRepo, "url = ?", t.Repo.URL)
		if result.Error != nil {
			errs.Panic(result.Error)
		}
		modRepo.Name = t.Repo.Name
		modRepo.URL = t.Repo.URL
		modRepo.Desc = t.Repo.Description

		orm.Model(&modRepo).Updates(&modRepo)
	}

	modTemplate := model.Template{
		Name: t.Name,
		URL:  t.URL,
		Desc: t.Description,
		Repo: []model.Repo{modRepo},
	}
	result = orm.Clauses(clause.Insert{Modifier: "OR IGNORE"}).Create(&modTemplate)
	if result.RowsAffected != 1 {
		result = orm.First(&modTemplate, "url = ?", t.URL)
		if result.Error != nil {
			errs.Panic(result.Error)
		}
		modTemplate.Name = t.Name
		modTemplate.URL = t.URL
		modTemplate.Desc = t.Description
		modTemplate.Repo = append(modTemplate.Repo, modRepo)

		orm.Model(&modTemplate).Updates(&modTemplate)
	}
}

// Repo is the repo struct
type Repo struct {
	Name        string `json:"name" yaml:"name"`
	URL         string `json:"url" yaml:"url"`
	Description string `json:"description" yaml:"description"`
}

// Load loads from a json or yaml file, depending on the file suffix.
func (m *Repo) Load(path string) error {
	return marshal.FromFile(m, path)
}

// Save saves to json or yaml file, depending on the file suffix.
func (m *Repo) Save(path string) error {
	return marshal.ToFile(m, path)
}

// Template is a template creator
type Template struct {
	Name        string `json:"name" yaml:"name"`
	URL         string `json:"url" yaml:"url"`
	Description string `json:"description" yaml:"description"`
	Repo        Repo
}

// Load loads from a json or yaml file
func (g *Template) Load(path string) error {
	return marshal.FromFile(g, path)
}

// Save saves to json or yaml file.
func (g *Template) Save(path string) error {
	return marshal.ToFile(g, path)
}
