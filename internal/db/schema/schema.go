package schema

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/crosseyed/prjstart/internal/config"
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/crosseyed/prjstart/internal/utils/dfaults"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	_ "github.com/mattn/go-sqlite3"
)

var DBFILE = "metadata.db"

// FKON - Query to turn on Foreign Key constraints
var FKON = "PRAGMA foreign_keys = ON"

var TBLMASTER = `
CREATE TABLE IF NOT EXISTS master (
	id integer not null primary key autoincrement,
	name text,
	url text,
	desc text
);
CREATE INDEX IF NOT EXISTS idx_master_name ON master(name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_master_url ON master(url);
`

var TBLTMPLS = `
CREATE TABLE IF NOT EXISTS templates (
	id integer not null primary key autoincrement,
	masterid integer, 
	name text, /* Suggested template name. This is set in the template metadata */
	url text,  /* URL to template */
	desc text, /* Description */
	FOREIGN KEY(masterid) REFERENCES master(id)
);
CREATE INDEX IF NOT EXISTS idx_templates_name ON templates (name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_templates_masterid_url ON templates(masterid, url);
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

// Schema handles creation of the databases
type Schema struct {
	db *sql.DB
}

// New generates a new schema
func New(driver, datasource string) *Schema {
	home := os.Getenv("HOME")
	confdir := config.USERCONFIGDIR
	fpath := filepath.Join(home, confdir, DBFILE)
	dbfile := dfaults.String(fpath, datasource)
	driver = dfaults.String("sqlite3", driver)
	utils.Touch(dbfile)
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on", dbfile)
	db, err := sql.Open("sqlite3", dsn)
	fmt.Println(dbfile)
	errutils.Epanicf("can not open file %s: %w", dbfile, err)
	return &Schema{
		// TODO: Depedency Injection
		db: db,
	}
}

// Create creates the schema
func (s *Schema) Create() {
	for _, query := range []string{FKON, TBLMASTER, TBLTMPLS, TBLTMPLSVER} {
		_, err := s.db.Exec(query)
		if err != nil {
			log.Fatalf("error creating database scheme: %v", err)
		}
	}
	return
}
