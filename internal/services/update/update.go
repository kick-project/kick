package update

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/apex/log"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/db"
	"github.com/kick-project/kick/internal/resources/gitclient"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/marshal"
	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
	"gorm.io/gorm"
)

// Update build metadata
type Update struct {
	ConfigFile  *config.File
	DB          *sql.DB
	ORM         *gorm.DB
	Log         *log.Logger
	MetadataDir string
}

// Build metadata. Conf defaults to globals.Config if Conf is nil.
func (m *Update) Build() error {
	conf := m.ConfigFile

	c := workers{
		wait: &sync.WaitGroup{},
		log:  m.Log,
	}

	churl := make(chan string, 64)
	chtemplates := make(chan *Template, 64)
	p := plumbing.New(m.MetadataDir)
	c.concurClones(6, p, churl, chtemplates)
	c.concurInserts(m.DB, chtemplates)

	for _, url := range conf.MasterURLs {
		c.wait.Add(1)
		churl <- url
	}

	// Wait for all all processing to finish
	c.wait.Wait()

	return nil
}

type workers struct {
	log  *log.Logger
	wait *sync.WaitGroup
}

// concurClones concurrent cloning of git repositories.
// where num is the number of concurrent downloads, churl is a string url and tchan is a channel of resulting templates.
func (c *workers) concurClones(num int, p *plumbing.Plumbing, churl <-chan string, tchan chan<- *Template) {
	for i := 0; i < num; i++ {
		go func() {
			for {
				url, ok := <-churl
				switch {
				case !ok:
					return
				default:
					c.processURL(url, p, tchan)
					c.wait.Done()
				}
			}
		}()
	}
}

func (c *workers) processURL(url string, p *plumbing.Plumbing, chtemplate chan<- *Template) {
	localpath, err := gitclient.Get(url, p)
	if errutils.Elogf("error: cloning repository: %w: skipping %s", err, url) {
		return
	}

	mpath := filepath.Clean(fmt.Sprintf("%s/master.yml", localpath))
	if errutils.Elogf("error: can not open %s: %w: skipping %s", mpath, err, url) {
		return
	}

	master := &Master{URL: url}
	err = master.Load(mpath)
	if errutils.Elogf("error: %w: skipping %s\n", err, url) {
		return
	}

	paths, err := filepath.Glob(filepath.Clean(fmt.Sprintf("%s/templates/*.yml", localpath)))
	if errutils.Elogf("error: getting a lists of paths: %w: skipping %s\n", err, url) {
		return
	}

	for _, curpath := range paths {
		t := &Template{}
		err := t.Load(curpath)
		if errutils.Elogf("error: loading template metadata from %s: %w: skipping", curpath, err) {
			continue
		}
		t.Master = *master
		c.wait.Add(1)
		chtemplate <- t
	}
}

// concurInserts populates the database
// where num is the number of concurrent routines
// and ch is the channel to read templates from.
func (c *workers) concurInserts(dbconn *sql.DB, ch <-chan *Template) {
	go func() {
		for {
			t, ok := <-ch
			switch {
			case !ok:
				return
			default:
				db.Lock()
				c.insert(dbconn, t)
				db.Unlock()
				c.wait.Done()
			}
		}
	}()
}

func (c *workers) insert(dbconn *sql.DB, t *Template) {
	insertMaster := `INSERT OR IGNORE INTO master (name, url, desc) VALUES (?, ?, ?)`
	c.log.Debugf(utils.SQL2fmt(insertMaster), t.Master.Name, t.Master.URL, t.Master.Description)
	_, err := dbconn.Exec(insertMaster, t.Master.Name, t.Master.URL, t.Master.Description)
	errutils.Efatalf("error: inserting master metadata: %w", err)

	insertTemplate := `INSERT OR IGNORE INTO templates (masterid, name, url, desc) SELECT master.id, ?, ?, ? FROM master WHERE master.url = ?`
	c.log.Debugf(utils.SQL2fmt(insertTemplate), t.Name, t.URL, t.Description, t.Master.URL)
	insertParams := []interface{}{t.Name, t.URL, t.Description, t.Master.URL}
	_, err = dbconn.Exec(insertTemplate, insertParams...)
	errutils.Efatalf("error: inserting template metadata: %w", err)
}

// Master is the master struct
type Master struct {
	Name        string `json:"name" yaml:"name"`
	URL         string `json:"url" yaml:"url"`
	Description string `json:"description" yaml:"description"`
}

// Load loads from a json or yaml file, depending on the file suffix.
func (m *Master) Load(path string) error {
	return marshal.UnmarshalFromFile(m, path)
}

// Save saves to json or yaml file, depending on the file suffix.
func (m *Master) Save(path string) error {
	return marshal.Marshal2File(m, path)
}

// Template is a template creator
type Template struct {
	Name        string `json:"name" yaml:"name"`
	URL         string `json:"url" yaml:"url"`
	Description string `json:"description" yaml:"description"`
	Master      Master
}

// Load loads from a json or yaml file
func (g *Template) Load(path string) error {
	return marshal.UnmarshalFromFile(g, path)
}

// Save saves to json or yaml file.
func (g *Template) Save(path string) error {
	return marshal.Marshal2File(g, path)
}
