package repo

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator"
	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di/callbacks"
	"github.com/kick-project/kick/internal/resources/client"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/logger"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/serialize"
)

// Repo build a repository repo
type Repo struct {
	client     *client.Client      // Git client
	makePlumb  callbacks.MakePlumb // Dependency Injector
	serialized serialize.RepoMain  // Serialized config
	errs       errs.HandlerIface   // Error handler
	log        logger.OutputIface  // Logger
}

// Options options for New
type Options struct {
	Client     *client.Client      `validate:"required"` // Git client
	MakePlumb  callbacks.MakePlumb `validate:"required"` // Dependency Injector
	ErrHandler errs.HandlerIface   `validate:"required"` // Error handler
	Log        logger.OutputIface  `validate:"required"` // Logger
}

// New construct a Repo object
func New(opts *Options) *Repo {
	err := validator.New().Struct(opts)
	if err != nil {
		panic(err)
	}
	r := &Repo{
		client:    opts.Client,
		makePlumb: opts.MakePlumb,
		errs:      opts.ErrHandler,
		log:       opts.Log,
	}
	return r
}

// Build build repo
func (m *Repo) Build() {
	m.load()
	m.download()
}

func (m *Repo) wd() string {
	wd, err := os.Getwd()
	m.errs.Panic(err)
	return wd
}

func (m *Repo) load() {
	fp := filepath.Join(m.wd(), "repo.yml")
	err := marshal.FromFile(&m.serialized, fp)
	m.errs.FatalF("Can not load file \"%s\": %v", fp, err)

	v := validator.New()
	err = v.Struct(&m.serialized)
	m.errs.FatalF("Can not load file \"%s\", invalid fields: %v", fp, err)
}

func (m *Repo) download() {
	destDir := filepath.Join(m.wd(), "templates")
	err := os.MkdirAll(destDir, 0755)
	errs.FatalF("Can create directory \"%s\": %v", destDir, err)

	v := validator.New()
	for _, url := range m.serialized.TemplateURLs {
		// Validate url
		err := v.Var(url, "url")
		if m.errs.LogF("Invalid url \"%s\": %v", url, err) {
			continue
		}

		// Get URL
		plumb := m.makePlumb(url, "")
		err = m.client.GetPlumb(plumb)
		if m.errs.LogF("Can not download \"%s\": %v", url, err) {
			continue
		}

		// Load .kick.yml
		var templateMain serialize.TemplateMain
		srcTemplate := filepath.Join(plumb.Path(), ".kick.yml")
		err = marshal.FromFile(&templateMain, srcTemplate)
		if m.errs.LogF("Can not load file \"%s\": %v", srcTemplate, err) {
			continue
		}

		// Validate .kick.yml
		err = v.Struct(&templateMain)
		if err != nil {
			var invalid []string
			for _, err := range err.(validator.ValidationErrors) {
				invalid = append(invalid, err.StructField())
			}
			m.log.Errorf("Can not load %s invalid fields: ", strings.Join(invalid, `,`))
			continue
		}

		// Copy object to "templates/*.yml" yaml file
		var templateElement serialize.RepoTemplateFile
		err = copier.Copy(&templateElement, &templateMain)
		if m.errs.LogF("Can not copy objects: %v", err) {
			continue
		}
		// Add URL
		templateElement.URL = url

		// Write "templates/*.yml" yaml file
		destRepoYAML := filepath.Join(destDir, templateElement.Name+".yml")
		err = marshal.ToFile(&templateElement, destRepoYAML)
		if m.errs.LogF("Can not save file \"%s\": %v", destRepoYAML, err) {
			continue
		}
	}
}
