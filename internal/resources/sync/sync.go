// Package sync synchronizes configuration between the database and
// corresponding file and synchronization of any downloads.
package sync

import (
	"fmt"
	"io"
	"time"

	"github.com/apex/log"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/db"
	"github.com/kick-project/kick/internal/resources/gitclient"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/utils/errutils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Sync synchronize database tables
type Sync struct {
	BasePath           string
	ORM                *gorm.DB           `copier:"must"`
	Config             *config.File       `copier:"must"`
	ConfigTemplatePath string             `copier:"must"`
	Log                *log.Logger        `copier:"must"`
	PlumbTemplates     *plumbing.Plumbing `copier:"must"`
	PlumbGlobal        *plumbing.Plumbing `copier:"must"`
	PlumbMaster        *plumbing.Plumbing `copier:"must"`
	Stderr             io.Writer          `copier:"must"`
	Stdout             io.Writer          `copier:"must"`
}

// Global syncs global data
func (s *Sync) Global() {
	rows, err := s.ORM.Model(&model.Global{}).Rows()
	errutils.Epanic(err)

	defer rows.Close()
	for rows.Next() {
		var global model.Global
		err := s.ORM.ScanRows(rows, &global)
		if err != nil {
			fmt.Fprintf(s.Stderr, "warning. can not scan table row from `global`: %v\n", err)
		}

		_, err = gitclient.Get(global.URL, s.PlumbGlobal)
		if err != nil {
			fmt.Fprintf(s.Stderr, "warning. can not download %s: %s\n", global.URL, err.Error())
			continue
		}
	}
}

// Master syncs master data
func (s *Sync) Master() {
}

// Templates synchronizes templates between the YAML configuration, database
// and its upstream version control repository.
func (s *Sync) Templates() {
	key := "installed"
	db.Lock()
	defer db.Unlock()
	// Reload configuration incase the file changed after creation of self.
	err := s.Config.Load()
	errutils.Epanic(err)
	t := time.Now()
	ts := t.Format("2006-01-02T15:04:05")
	for _, item := range s.Config.Templates {
		_, err := gitclient.Get(item.URL, s.PlumbTemplates)
		if err != nil {
			fmt.Fprintf(s.Stderr, "warning. can not download %s: %s\n", item.URL, err.Error())
		}
		inst := model.Installed{
			Handle:   item.Handle,
			Template: item.Template,
			Origin:   item.Origin,
			URL:      item.URL,
			Desc:     item.Desc,
			Time:     t,
		}
		result := s.ORM.Clauses(clause.Insert{Modifier: "OR REPLACE"}).Create(&inst)
		errutils.Epanic(result.Error)
		if result.RowsAffected != 1 {
			panic("failed to insert into 'installed' table")
		}
	}

	result := s.ORM.Raw(`DELETE FROM installed WHERE time < ?`, ts)
	errutils.Epanic(result.Error)

	syn := model.Sync{
		Key:        key,
		LastUpdate: t,
	}
	result = s.ORM.Clauses(clause.Insert{Modifier: "OR REPLACE"}).Create(&syn)
	errutils.Epanic(result.Error)
	if result.RowsAffected != 1 {
		panic("failed to insert into 'sync' table")
	}
}
