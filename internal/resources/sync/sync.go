// Package sync synchronizes configuration between the database and
// corresponding file and synchronization of any downloads.
package sync

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/client"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/logger"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/resources/model/clauses"
	"github.com/kick-project/kick/internal/resources/serialize"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO: Unit test coverage. Increase from 64%

// Sync synchronize database tables
//go:generate ifacemaker -f sync.go -s Sync -p sync -i SyncIface -o sync_interfaces.go -c "AUTO GENERATED. DO NOT EDIT."
type Sync struct {
	client             *client.Client
	orm                *gorm.DB
	config             *config.File
	configTemplatePath string
	log                logger.OutputIface
	stderr             io.Writer
	stdout             io.Writer
}

// Options options to constructor
type Options struct {
	Client             *client.Client     `validate:"required"`
	Config             *config.File       `validate:"required"`
	ConfigTemplatePath string             `validate:"required"`
	Log                logger.OutputIface `validate:"required"`
	ORM                *gorm.DB           `validate:"required"`
	Stderr             io.Writer          `validate:"required"`
	Stdout             io.Writer          `validate:"required"`
}

// New construct Sync object
func New(opts *Options) *Sync {
	return &Sync{
		client:             opts.Client,
		config:             opts.Config,
		configTemplatePath: opts.ConfigTemplatePath,
		log:                opts.Log,
		orm:                opts.ORM,
		stderr:             opts.Stderr,
		stdout:             opts.Stdout,
	}

}

// Repo syncs repo data
func (s *Sync) Repo() {
	repos := s.processRepos()
	s.processTemplates(repos)
}

func (s *Sync) processRepos() (repos []*model.Repo) {
	rows, err := s.orm.Model(&model.Repo{}).Rows()
	errs.Panic(err)

	defer rows.Close()
	for rows.Next() {
		var repo model.Repo
		err := s.orm.ScanRows(rows, &repo)
		if err != nil {
			fmt.Fprintf(s.stderr, "warning. can not scan table row from `repo`: %v\n", err)
		}

		if repo.URL == "none" {
			continue
		}

		repos = append(repos, &repo)
	}
	return
}

func (s *Sync) processTemplates(repos []*model.Repo) {
	for _, repo := range repos {
		path, err := s.downloadRepo(repo.URL)
		if err != nil {
			continue
		}

		repoPath := filepath.Clean(fmt.Sprintf("%s/%s", path, "repo.yml"))
		repoSerialize, err := s.loadRepo(repoPath)
		if err != nil {
			continue
		}

		err = copier.Copy(&repo, &repoSerialize)
		if errs.LogF("Can not copy object: %v", err) {
			continue
		}
		result := s.orm.Clauses(clauses.OrIgnore).Create(&repo)
		if errs.LogF("Can not insert into repo: %v", result.Error) {
			continue
		} else if result.RowsAffected == 0 {
			result2 := s.orm.Model(&model.Repo{}).Updates(&repo)
			if errs.LogF("Can not update repo table: %v", result2.Error) {
				continue
			} else if result2.RowsAffected == 0 {
				errs.LogF("%v", fmt.Errorf("Can not update repo table"))
				continue
			}
		}

		s.loadTemplates(repo, filepath.Join(path, "templates"))
	}
}

// downloadRepo downloads repo repo
func (s *Sync) downloadRepo(url string) (path string, err error) {
	p, err := s.client.GetRepo(url, "")
	if errs.LogF("warning. can not download %s: %v\n", url, err) {
		return
	}
	path = p.Path()
	return
}

// loadRepo loads from a repo YAML file
func (s *Sync) loadRepo(path string) (repo serialize.RepoMain, err error) {
	err = marshal.FromFile(&repo, path)
	if errs.LogF("warning. unable to unmarshal file \"%s\": %v", path, err) {
		return
	}
	return
}

// TODO: Add unit test for loadTemplates

// loadTemplates loads templates from a repo file
func (s *Sync) loadTemplates(repo *model.Repo, templatedir string) {
	matches, err := filepath.Glob(templatedir + "/*.yml")
	if errs.LogF("Can not load templates from \"%s\": %v", templatedir, err) {
		return
	}
	for _, match := range matches {
		var (
			templateElement serialize.RepoTemplateFile
			templateModel   model.Template
		)

		err := marshal.FromFile(&templateElement, match)
		if errs.LogF("Can not load template file \"%s\": %v", match, err) {
			continue
		}

		err = copier.Copy(&templateModel, &templateElement)
		if errs.LogF("Can not copy object: %v", err) {
			continue
		}

		templateModel.Repo = append(templateModel.Repo, *repo)

		result := s.orm.Clauses(clauses.OrReplace).Create(&templateModel)
		if errs.LogF("Can not load template file \"%s\" into database: %v", match, result.Error) {
			continue
		}
	}
}

// Files synchronizes templates between the YAML configuration, database
// and its upstream version control repository.
func (s *Sync) Files() {
	key := "installed"
	// Reload configuration incase the file changed after creation of self.
	err := s.config.Load()
	errs.Panic(err)
	t := time.Now()
	ts := t.Format("2006-01-02T15:04:05")
	for _, item := range s.config.Templates {
		_, err := s.client.GetTemplate(item.URL, "")
		if err != nil {
			fmt.Fprintf(s.stderr, "warning. can not download %s: %s\n", item.URL, err.Error())
		}
		inst := model.Installed{
			Handle:   item.Handle,
			Template: item.Template,
			Origin:   item.Origin,
			URL:      item.URL,
			Desc:     item.Desc,
			Time:     t,
		}
		result := s.orm.Clauses(clause.Insert{Modifier: "OR REPLACE"}).Create(&inst)
		errs.Panic(result.Error)
		if result.RowsAffected != 1 {
			panic("failed to insert into 'installed' table")
		}
	}

	result := s.orm.Raw(`DELETE FROM installed WHERE time < ?`, ts)
	errs.Panic(result.Error)

	syn := model.Sync{
		Key:        key,
		LastUpdate: t,
	}
	result = s.orm.Clauses(clause.Insert{Modifier: "OR REPLACE"}).Create(&syn)
	errs.Panic(result.Error)
	if result.RowsAffected != 1 {
		panic("failed to insert into 'sync' table")
	}
}
