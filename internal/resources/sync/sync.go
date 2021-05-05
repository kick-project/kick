// Package sync synchronizes configuration between the database and
// corresponding file and synchronization of any downloads.
package sync

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/apex/log"
	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/gitclient"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/resources/serialize"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/marshal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Sync synchronize database tables
type Sync struct {
	ORM                *gorm.DB           `copier:"must"`
	Config             *config.File       `copier:"must"`
	ConfigTemplatePath string             `copier:"must"`
	Log                *log.Logger        `copier:"must"`
	PlumbTemplates     *plumbing.Plumbing `copier:"must"`
	PlumbMaster        *plumbing.Plumbing `copier:"must"`
	Stderr             io.Writer          `copier:"must"`
	Stdout             io.Writer          `copier:"must"`
}

// Master syncs master data
func (s *Sync) Master() {
	rows, err := s.ORM.Model(&model.Master{}).Rows()
	errutils.Epanic(err)

	defer rows.Close()
	for rows.Next() {
		var master model.Master
		err := s.ORM.ScanRows(rows, &master)
		if err != nil {
			fmt.Fprintf(s.Stderr, "warning. can not scan table row from `global`: %v\n", err)
		}

		if master.URL == "none" {
			continue
		}

		path, err := s.downloadMaster(master.URL)
		if err != nil {
			continue
		}

		masterPath := filepath.Clean(fmt.Sprintf("%s/%s", path, "master.yml"))
		masterSerialize, err := s.loadMaster(masterPath)
		if err != nil {
			continue
		}
		_ = copier.Copy(&master, &masterSerialize)
		s.ORM.Model(&model.Master{}).Updates(&master)
	}
}

// downloadMaster downloads master repo
func (s *Sync) downloadMaster(url string) (path string, err error) {
	path, err = gitclient.Get(url, s.PlumbMaster)
	if errutils.Elogf("warning. can not download %s: %v\n", url, err) {
		return
	}

	return
}

// loadMaster loads global file
func (s *Sync) loadMaster(path string) (master serialize.Master, err error) {
	err = marshal.UnmarshalFromFile(&master, path)
	if errutils.Elogf("warning. unable to unmarshal file \"%s\": %v", path, err) {
		return
	}
	return
}

// Files synchronizes templates between the YAML configuration, database
// and its upstream version control repository.
func (s *Sync) Files() {
	key := "installed"
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
