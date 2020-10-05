package metadata

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/crosseyed/prjstart/internal/config"
	"github.com/crosseyed/prjstart/internal/db"
	"github.com/crosseyed/prjstart/internal/db/schema"
	"github.com/crosseyed/prjstart/internal/gitclient"
	"github.com/crosseyed/prjstart/internal/gitclient/plumbing"
	"github.com/crosseyed/prjstart/internal/globals"
	"github.com/crosseyed/prjstart/internal/marshal"
	"github.com/crosseyed/prjstart/internal/utils/dfaults"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
)

// Build metadata. Conf defaults to globals.Config if Conf is nil.
func Build(conf *config.Config) error {
	// TODO: Inject dependency
	v := dfaults.Interface(&globals.Config, conf)
	conf, ok := v.(*config.Config)
	if !ok {
		return errors.New("error: can not build. config is an unknown data type")
	}

	// TODO: Inject dependency
	path := filepath.Join(conf.Home, ".prjstart", "metadata", "metadata.db")
	dirname := filepath.Dir(path)
	err := os.MkdirAll(dirname, 0755)
	errutils.Efatalf("error: %w", err)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		s := schema.New("sqlite3", path)
		s.Create()
	} else {
		errutils.Efatalf("error: %w", err)
	}

	c := workers{
		wait: &sync.WaitGroup{},
	}

	churl := make(chan string, 64)
	chtemplates := make(chan *Template, 64)
	metadatadir := filepath.Join(conf.Home, ".prjstart", "metadata")
	p := plumbing.New(metadatadir)
	c.runClones(6, p, churl, chtemplates)
	c.runInserts(fmt.Sprintf("file:%s?_foreign_keys=on", path), chtemplates)

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

// runClones concurrent cloning of git repositories.
// where num is the number of concurrent downloads, churl is a string url and tchan is a channel of resulting templates.
func (c *workers) runClones(num int, p *plumbing.Plumbing, churl <-chan string, tchan chan<- *Template) {
	for i := 0; i < num; i++ {
		go func() {
			for {
				url, ok := <-churl
				switch {
				case !ok:
					return
				default:
					c.processUrl(url, p, tchan)
					c.wait.Done()
				}
			}
		}()
	}
}

func (c *workers) processUrl(url string, p *plumbing.Plumbing, chtemplate chan<- *Template) {
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
		fmt.Fprintln(os.Stderr, t)
		fmt.Fprintln(os.Stderr, master)
		c.wait.Add(1)
		chtemplate <- t
	}
}

// runInserts populates the database
// where num is the number of concurrent routines
// and ch is the channel to read templates from.
func (c *workers) runInserts(dsn string, ch <-chan *Template) {
	// TODO: I hate my UGGGGLLLLLY code. Needs dependency injection
	dbconn, err := sql.Open("sqlite3", dsn)
	errutils.Epanicf("error: could not connect to the metadata database: %w", err)
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
	// TODO: DI
	log.Printf(strings.ReplaceAll(insertTemplate, "?", `"%s"`), insertParams...)
	log.Println()
	_, err = dbconn.Exec(insertTemplate, insertParams...)
	errutils.Efatalf("error: inserting template metadata: %w", err)
}

// Master is the master struct
type Master struct {
	Name        string `json:"name" yaml:"name"`
	URL         string
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
