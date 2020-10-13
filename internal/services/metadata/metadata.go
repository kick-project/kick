package metadata

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/resources/db"
	"github.com/crosseyed/prjstart/internal/resources/gitclient"
	"github.com/crosseyed/prjstart/internal/resources/gitclient/plumbing"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/crosseyed/prjstart/internal/utils/marshal"
	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
)

// Metadata build metadata
type Metadata struct {
	conf         *config.File
	metadatapath string
	db           *sql.DB
}

// Options options for metadata.New
type Options struct {
	ConfigFile  *config.File // Contents of the main config file in code
	MetadataDir string
	DB          *sql.DB
}

// New create an instance of Metadata.
// Panics if opts.ConfigFile is nil or DBPath is an empty string.
func New(opts Options) *Metadata {
	if opts.ConfigFile == nil {
		panic("opts.ConfigFile can not be nil")
	}
	if opts.MetadataDir == "" {
		panic("opts.MetadataPath can not be an empty string")
	}
	m := &Metadata{
		conf:         opts.ConfigFile,
		metadatapath: opts.MetadataDir,
		db:           opts.DB,
	}
	return m
}

// Build metadata. Conf defaults to globals.Config if Conf is nil.
func (m *Metadata) Build() error {
	conf := m.conf

	c := workers{
		wait: &sync.WaitGroup{},
	}

	churl := make(chan string, 64)
	chtemplates := make(chan *Template, 64)
	p := plumbing.New(m.metadatapath)
	c.concurClones(6, p, churl, chtemplates)
	c.concurInserts(m.db, chtemplates)

	for _, url := range conf.MasterURLs {
		c.wait.Add(1)
		churl <- url
	}

	// Wait for all all processing to finish
	c.wait.Wait()

	return nil
}

type workers struct {
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
	_, err := dbconn.Exec(insertMaster, t.Master.Name, t.Master.URL, t.Master.Description)
	errutils.Efatalf("error: inserting master metadata: %w", err)

	insertTemplate := `INSERT OR IGNORE INTO templates (masterid, name, url, desc) SELECT master.id, ?, ?, ? FROM master WHERE master.url = ?`
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
	return marshal.UnmarshalFile(m, path)
}

// Save saves to json or yaml file, depending on the file suffix.
func (m *Master) Save(path string) error {
	return marshal.MarshalFile(m, path)
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
	return marshal.UnmarshalFile(g, path)
}

// Save saves to json or yaml file.
func (g *Template) Save(path string) error {
	return marshal.MarshalFile(g, path)
}
