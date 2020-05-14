package dbinit

import (
	"log"
	"os"
	"path/filepath"

	"github.com/crosseyed/prjstart/internal/config"
	"github.com/crosseyed/prjstart/internal/db"
	"github.com/crosseyed/prjstart/internal/utils/dfaults"
	_ "github.com/mattn/go-sqlite3"
)

var DBFILE = "metadata.db"

var TBLORG = `
CREATE TABLE IF NOT EXISTS global (
	id integer not null primary key autoincrement,
	url text,
	name text,
	desc text
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_global_url ON global (url);
CREATE INDEX IF NOT EXISTS idx_global_name ON global (name);
`

var TBLMASTER = `
CREATE TABLE IF NOT EXISTS master (
	id integer not null primary key autoincrement,
	globalid integer,
	name text,
	url text,
	desc text,
	FOREIGN KEY(globalid) REFERENCES global(id)
);
CREATE INDEX IF NOT EXISTS idx_master_name ON master (name);
CREATE INDEX IF NOT EXISTS idx_master_global_fk ON master (globalid);
CREATE UNIQUE INDEX IF NOT EXISTS idx_master_url ON master (url);
`

var TBLTMPLS = `
CREATE TABLE IF NOT EXISTS templates (
	id integer not null primary key autoincrement,
	masterid integer,
	name text,
	url text,
	desc text,
	installed bool,
	FOREIGN KEY(masterid) REFERENCES master(id)
);
CREATE INDEX IF NOT EXISTS idx_templates_masterid_name_url ON templates (masterid, name, url);
`

var TBLTMPLSVER = `
CREATE TABLE IF NOT EXISTS versions (
	id integer not null primary key autoincrement,
	version text,
	templatesid integer not null,
	FOREIGN KEY(templatesid) REFERENCES templates(id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_versions_templateid_version ON versions (templatesid, version);
`

var INSERTGLOBAL string = `INSERT OR REPLACE INTO global (url, name, desc) VALUES (?, ?, ?)`
var INSERTMASTER string = `INSERT OR REPLACE INTO master (url, name, desc, globalid) SELECT ?, ?, ?, id FROM global WHERE url=?`

type Init struct {
	DB *db.DB
}

func New(datasource string) *Init {
	home := os.Getenv("HOME")
	confdir := config.USERCONFIGDIR
	fpath := filepath.Join(home, confdir, DBFILE)
	dbfile := dfaults.String(fpath, datasource)
	return &Init{
		DB: db.New("sqlite3", dbfile),
	}
}

func (s *Init) Init() {
	s.DB.Lock()
	s.DB.Open()
	defer s.DB.Unlock()
	defer s.DB.Open()
	s.createSchema()
	return
}

func (s *Init) createSchema() {
	for _, query := range []string{TBLORG, TBLMASTER, TBLTMPLS, TBLTMPLSVER} {
		s.execWrapper(query)
	}
	return
}

func (s *Init) execWrapper(query string, queryArgs ...interface{}) {
	_, err := s.DB.Exec(query, queryArgs...)
	if err != nil {
		log.Fatalf("SQL Error creating defaults: %v", err)
	}
}
