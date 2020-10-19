// Package sync synchronizes configuration between the database and
// corresponding file and synchronization of any downloads.
package sync

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/apex/log"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/db"
	"github.com/kick-project/kick/internal/resources/gitclient"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/utils/errutils"
)

//
// SQL
//

var selectSync = `SELECT
  count(CASE WHEN lastupdate < ? THEN 1 ELSE NULL END) AS status, 
  count(*) AS haskey FROM sync WHERE key=?`

var insertReplaceSync = `INSERT OR REPLACE INTO sync (key, lastupdate) VALUES (?, ?)`

var insertReplaceInstalled = `INSERT OR REPLACE INTO installed (handle, template, origin, url, desc, time) VALUES (?, ?, ?, ?, ?, ?)`

var deleteMissing = `DELETE FROM installed WHERE time < ?`

//
// Go
//

// Sync synchronize database tables
type Sync struct {
	DB                 *sql.DB
	Config             *config.File
	ConfigTemplatePath string
	Log                *log.Logger
	Plumb              *plumbing.Plumbing
	Stderr             io.Writer
	Stdout             io.Writer
}

// Check returns true if a file needs synchronizing.
// key is the key the database which holds the last update time and file is the path to stat
func (s *Sync) Check(key, file string) bool {
	file, err := filepath.Abs(filepath.Clean(file))
	errutils.Epanicf("%w", err)
	info, err := os.Stat(file)
	errutils.Epanicf("%w", err)
	ts := info.ModTime().Format("2006-01-02T15:04:05")
	row := s.DB.QueryRow(selectSync, ts, key)
	update := 0
	haskey := 0
	err = row.Scan(&update, &haskey)
	errutils.Epanicf("%w", err)
	return update == 1 || haskey == 0
}

// Templates synchronizes templates between the YAML configuration, database
// and its upstream version control repository.
func (s *Sync) Templates() {
	key := "installed"
	if !s.Check("installed", s.ConfigTemplatePath) {
		return
	}
	db.Lock()
	defer db.Unlock()
	// Reload configuration incase the file changed after creation of self.
	err := s.Config.Load()
	errutils.Epanic(err)
	t := time.Now()
	ts := t.Format("2006-01-02T15:04:05")
	for _, item := range s.Config.Templates {
		_, err := gitclient.Get(item.URL, s.Plumb)
		if err != nil {
			fmt.Fprintf(s.Stderr, "warning. can not download %s: %s\n", item.URL, err.Error())
		}
		_, err = s.DB.Exec(insertReplaceInstalled, item.Handle, item.Template, item.Origin, item.URL, item.Desc, ts)
		errutils.Epanicf("%w", err)
	}

	_, err = s.DB.Exec(deleteMissing, ts)
	errutils.Epanicf("%w", err)

	_, err = s.DB.Exec(insertReplaceSync, key, ts)
	errutils.Epanicf("%w", err)
}
