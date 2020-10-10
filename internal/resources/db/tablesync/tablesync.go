// Package tablesync synchronizes database tables with the coresponding file.
package tablesync

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/resources/db"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

//
// SQL
//

var selectSync = `SELECT
  count(CASE WHEN lastupdate < ? THEN 1 ELSE NULL END) AS status, 
  count(*) AS haskey FROM sync WHERE key=?`

var insertReplaceSync = `INSERT OR REPLACE INTO sync (key, lastupdate) VALUES (?, ?)`

var insertReplaceInstalled = `INSERT OR IGNORE INTO installed (name, url, time) VALUES (?, ?, ?)`

var deleteMissing = `DELETE FROM installed WHERE time < ?`

//
// Go
//

// TableSync synchronize database tables
type TableSync struct {
	db         *sql.DB
	conf       *config.File
	configPath string
}

// Options options to New
type Options struct {
	DB         *sql.DB
	Config     *config.File
	ConfigPath string
}

// New returns a *TableSync object
func New(opts Options) *TableSync {
	if opts.DB == nil {
		panic("opts.DB is nil")
	}
	if opts.Config == nil {
		panic("opts.Config is nil")
	}
	s := &TableSync{
		db:   opts.DB,
		conf: opts.Config,
	}
	return s
}

// Check returns true if a file needs synchronizing.
// key is the key the database which holds the last update time and file is the path to stat
func (s *TableSync) Check(key, file string) bool {
	file, err := filepath.Abs(filepath.Clean(file))
	errutils.Epanicf("%w", err)
	info, err := os.Stat(file)
	errutils.Epanicf("%w", err)
	ts := info.ModTime().Format("2006-01-02T15:04:05")
	row := s.db.QueryRow(selectSync, ts, key)
	update := 0
	haskey := 0
	err = row.Scan(&update, &haskey)
	errutils.Epanicf("%w", err)
	return update == 1 || haskey == 0
}

// SyncInstalled syncs database/configuration file for installed components.
func (s *TableSync) SyncInstalled() {
	key := "installed"
	if !s.Check("installed", s.configPath) {
		return
	}
	t := time.Now()
	ts := t.Format("2006-01-02T15:04:05")
	db.Lock()
	defer db.Unlock()
	for _, item := range s.conf.TemplateURLs {
		_, err := s.db.Exec(insertReplaceInstalled, item.Name, item.URL, ts)
		errutils.Epanicf("%w", err)
	}

	_, err := s.db.Exec(deleteMissing, ts)
	errutils.Epanicf("%w", err)

	_, err = s.db.Exec(insertReplaceSync, key, ts)
	errutils.Epanicf("%w", err)
}
