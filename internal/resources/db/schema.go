package db

import (
	"database/sql"

	"github.com/kick-project/kick/internal/utils/errutils"
)

var tblMaster = `
CREATE TABLE IF NOT EXISTS master (
	id integer not null primary key autoincrement,
	name text,
	url text,
	desc text
);
CREATE INDEX IF NOT EXISTS idx_master_name ON master(name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_master_url ON master(url);
INSERT OR IGNORE INTO master (name, url, desc) VALUES ("local", "none", "This template is generated locally")
`

var tblTemplate = `
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

var tblInstalled = `
CREATE TABLE IF NOT EXISTS installed (
	id integer not null primary key autoincrement,
	handle text,
	template text,
	origin text,
	url text,
	vcsref text, 
	desc text,
	time text
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_installed_handle ON installed(handle);
CREATE UNIQUE INDEX IF NOT EXISTS idx_installed_handle_origin_url ON installed(handle, origin, url);
`

var tblSync = `	
CREATE TABLE IF NOT EXISTS sync (
	key text,
	lastupdate text /* datetime as text */
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sync_key ON sync(key);
CREATE INDEX IF NOT EXISTS idx_sync_lastupdate ON sync(lastupdate);
`

var tblVersions = `
CREATE TABLE IF NOT EXISTS versions (
	id integer not null primary key autoincrement,
	version text,
	templatesid integer not null,
	FOREIGN KEY(templatesid) REFERENCES templates(id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_versions_templateid_version ON versions (templatesid, version);
`

// CreateSchema creates a new schema
func CreateSchema(dbconn *sql.DB) {
	Lock()
	defer Unlock()
	for _, query := range []string{tblMaster, tblTemplate, tblSync, tblInstalled, tblVersions} {
		_, err := dbconn.Exec(query)
		if err != nil {
			errutils.Efatalf("error creating database scheme: %v", err)
		}
	}
}
