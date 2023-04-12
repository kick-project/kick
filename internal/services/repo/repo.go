package repo

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/go-playground/validator"
	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/client"
	"github.com/kick-project/kick/internal/resources/client/plumb"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/config/configtemplate"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/logger"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/resources/serialize"
	"github.com/kick-project/kick/internal/resources/vcs"
	"github.com/olekukonko/tablewriter"
	"gorm.io/gorm"
)

// Repo build a repository repo
//
//go:generate ifacemaker -f repo.go -s Repo -p repo -i RepoIface -o repo_interfaces.go -c "AUTO GENERATED. DO NOT EDIT."
type Repo struct {
	client     *client.Client      // Git client
	conf       *config.File        // Config file
	serialized serialize.RepoMain  // Serialized config
	errs       errs.HandlerIface   // Error handler
	log        logger.OutputIface  // Logger
	orm        *gorm.DB            // GoRM
	stdout     io.Writer           // Stdout
	twriter    *tablewriter.Table  // Table writer
	valid      *validator.Validate // Validation
	vcs        *vcs.VCS            // Version control repo
}

// Options options for New
type Options struct {
	Client      *client.Client      `validate:"required"` // Git client
	Conf        *config.File        `validate:"required"` // Config file
	ErrHandler  errs.HandlerIface   `validate:"required"` // Error handler
	Log         logger.OutputIface  `validate:"required"` // Logger
	ORM         *gorm.DB            `validate:"required"` // GORM
	Stdout      io.Writer           `validate:"required"` // Writer
	TableWriter *tablewriter.Table  `validate:"required"` // Tablewriter
	Valid       *validator.Validate `validate:"required"` // Validator
	VCS         *vcs.VCS            `validate:"required"` // Version Control Repo
}

// New construct a Repo object
func New(opts *Options) *Repo {
	r := &Repo{
		client:  opts.Client,
		conf:    opts.Conf,
		errs:    opts.ErrHandler,
		log:     opts.Log,
		orm:     opts.ORM,
		stdout:  opts.Stdout,
		twriter: opts.TableWriter,
		valid:   opts.Valid,
		vcs:     opts.VCS,
	}
	return r
}

// Build build repo
func (r *Repo) Build() {
	r.loadRepo()
	r.buildRepo()
}

func (r *Repo) wd() string {
	wd, err := os.Getwd()
	r.errs.Panic(err)
	return wd
}

func (r *Repo) loadRepo() {
	fp := filepath.Join(r.wd(), "repo.yml")
	err := marshal.FromFile(&r.serialized, fp)
	r.errs.FatalF("Can not load file \"%s\": %v", fp, err)

	v := validator.New()
	err = v.Struct(&r.serialized)
	r.errs.FatalF("Can not load file \"%s\", invalid fields: %v", fp, err)
}

func (r *Repo) buildRepo() {
	destDir := filepath.Join(r.wd(), "templates")
	err := os.MkdirAll(destDir, 0755)
	errs.FatalF("Can create directory \"%s\": %v", destDir, err)

	for _, url := range r.serialized.TemplateURLs {
		plu, ok := r.downloadTemplate(url)
		if !ok {
			continue
		}
		r.constructRepo(destDir, plu)
	}
}

func (r *Repo) downloadTemplate(url string) (plu *plumb.Plumb, ok bool) {
	// Validate url
	err := r.valid.Var(url, "url")
	if r.errs.LogF("Invalid url \"%s\": %v", url, err) {
		return nil, false
	}

	// Get URL
	plumb, err := r.client.GetTemplate(url, "")
	if r.errs.LogF("Can not download \"%s\": %v", url, err) {
		return nil, false
	}
	return plumb, true
}

func (r *Repo) constructRepo(destDir string, plu *plumb.Plumb) bool {
	// Load .kick.yml
	var templateMain configtemplate.TemplateMain
	srcTemplate := filepath.Join(plu.Path(), ".kick.yml")
	err := marshal.FromFile(&templateMain, srcTemplate)
	if r.errs.LogF("Can not load file \"%s\": %v", srcTemplate, err) {
		return false
	}

	// Validate .kick.yml
	err = r.valid.Struct(&templateMain)
	if err != nil {
		var invalid []string
		for _, err := range err.(validator.ValidationErrors) {
			invalid = append(invalid, err.StructField())
		}
		r.log.Errorf("Can not load %s invalid fields: ", strings.Join(invalid, `,`))
		return false
	}

	// Copy object to "templates/*.yml" yaml file
	var templateElement serialize.RepoTemplateFile
	err = copier.Copy(&templateElement, &templateMain)
	if r.errs.LogF("Can not copy objects: %v", err) {
		return false
	}
	// Add URL
	templateElement.URL = plu.URL()

	// Add Version
	templateElement.Versions = r.versions(plu)

	// Write "templates/*.yml" yaml file
	destRepoYAML := filepath.Join(destDir, templateElement.Name+".yml")
	err = marshal.ToFile(&templateElement, destRepoYAML)
	if r.errs.LogF("Can not save file \"%s\": %v", destRepoYAML, err) { // nolint
		return false
	}
	return true
}

func (r *Repo) versions(plu *plumb.Plumb) []string {
	versStr := []string{}
	repo, err := r.vcs.Open(plu.Path())
	r.errs.FatalF(`error opening %s: %w`, plu.Path(), err)
	// Sort verions
	var versions semver.Versions
	for _, v := range repo.Versions() {
		curver := semver.New(v)
		versions = append(versions, curver)
		_ = curver
	}
	sort.Sort(&versions)
	for _, v := range versions {
		versStr = append(versStr, v.String())
	}
	return versStr
}

// List list repositories
func (r *Repo) List() {
	result := r.orm.Model(&model.Repo{}).Select(`name, url`).Where(`name != ?`, "local")
	r.errs.FatalF(`database query error: %w`, result.Error)
	r.twriter.SetAlignment(tablewriter.ALIGN_LEFT)
	r.twriter.SetHeader([]string{"repo", "url"})
	rows, err := result.Rows()
	r.errs.FatalF(`database query error: %w`, err)
	for rows.Next() {
		var (
			name sql.NullString
			url  sql.NullString
		)
		err = rows.Scan(&name, &url)
		r.errs.FatalF(`table scan error: %w`, err)
		r.twriter.Append([]string{name.String, url.String})
	}
	r.twriter.Render()
}

// Info information on repositories
func (r *Repo) Info(repo string) {
	repoModel := &model.Repo{}
	_ = r.orm.First(repoModel, `name != ? and name = ?`, "local", repo)
	fmt.Fprintf(r.stdout, `name: %s
url: %s
description: %s`, repoModel.Name, repoModel.URL, repoModel.Desc)
}
