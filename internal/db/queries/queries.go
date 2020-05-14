package queries

import (
	"database/sql"

	"github.com/crosseyed/prjstart/internal/db"
	"github.com/crosseyed/prjstart/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

/*
 * SQL Queries
 */

var (
	// Populate global
	INSERTGLOBAL string = `INSERT OR REPLACE INTO global (url, name, desc) VALUES (?, ?, ?);`
	// Populate master
	INSERTMASTERBYURL string = `INSERT OR REPLACE INTO master (globalid, url, name, desc) SELECT id, "?", "?", "?" FROM global WHERE url=?;`
	// Populate template
	INSERTTEMPLATEBYURL string = `INSERT OR REPLACE INTO templates (masterid, url, name, desc) SELECT id, "?", "?", "?" FROM master WHERE url=?;`

	// Search templates
	MATCHTEMPLATE string = `
SELECT
	templates.name AS name,
	templates.url AS url,
	templates.desc AS desc,
	master.name AS master_name,
	master.url AS master_url,
	master.desc AS master_desc,
	global.name AS global_name,
	global.url AS global_url,
	global.desc AS global_desc,
	COALESCE(installed, false) AS installed FROM templates
LEFT JOIN master
LEFT JOIN global
WHERE templates.name = ?
UNION
SELECT
	templates.name AS name,
	templates.url AS url,
	templates.desc AS desc,
	master.name AS master_name,
	master.url AS master_url,
	master.desc AS master_desc,
	global.name AS global_name,
	global.url AS global_url,
	global.desc AS global_desc,
	COALESCE(installed, false) AS installed FROM templates
LEFT JOIN master
LEFT JOIN global
WHERE templates.name!=? AND templates.name LIKE '%?%'
ORDER BY templates.name, global.name, master.name
`
)

type SQL struct {
	DB *db.DB
}

func New(driver, datasource string) *SQL {
	return &SQL{
		DB: db.New(driver, datasource),
	}
}

func (s *SQL) InsertGlobal(url, name, desc string) (sql.Result, error) {
	s.DB.Lock()
	s.DB.Open()
	defer s.DB.Unlock()
	defer s.DB.Close()
	return s.DB.Exec(INSERTGLOBAL, url, name, desc)
}

func (s *SQL) InsertMasterByURL(globalurl, url, name, desc string) (sql.Result, error) {
	s.DB.Lock()
	s.DB.Open()
	defer s.DB.Unlock()
	defer s.DB.Close()
	return s.DB.Exec(INSERTMASTERBYURL, url, name, desc, globalurl)
}

func (s *SQL) InsertTemplateByURL(masterurl, url, name, desc string) (sql.Result, error) {
	s.DB.Lock()
	s.DB.Open()
	defer s.DB.Unlock()
	defer s.DB.Close()
	return s.DB.Exec(INSERTTEMPLATEBYURL, url, name, desc, masterurl)
}

func (s *SQL) SearchTemplate(template string, fnrow func(name, url, desc string, master_name, master_url, master_desc, global_name, global_url, global_desc string)) {
	s.DB.Lock()
	db := s.DB.Open()
	defer s.DB.Unlock()
	defer s.DB.Close()

	rows, err := db.Query(MATCHTEMPLATE, template, template, template)
	utils.ChkErr(err, utils.Efatalf, "Can not Query database: %v", err)

	for rows.Next() {
		var (
			name        string
			url         string
			desc        string
			master_name string
			master_url  string
			master_desc string
			global_name string
			global_url  string
			global_desc string
			installed   bool
		)
		err := rows.Scan(&name, &url, &desc, &master_name, &master_url, &master_desc, &global_name, &global_url, &global_desc, &installed)
		utils.ChkErr(err, utils.Efatalf, "Can not scan row: %v", err)

		fnrow(name, url, desc, master_name, master_url, master_desc, global_name, global_url, global_desc)
	}
}
