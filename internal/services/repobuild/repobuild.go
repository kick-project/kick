package repobuild

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator"
	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/gitclient"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/serialize"
)

// RepoBuild build a repository repo
type RepoBuild struct {
	WD         string              // Working Directory
	Plumb      plumbing.Plumbing   // Plumbing object
	Serialized serialize.RepoMain  // Serialized config
	Validate   *validator.Validate // Validation
	ErrHandler *errs.Errors        // Error handler
	Log        *log.Logger         // Logger
}

// Make build repo
func (m *RepoBuild) Make() {
	m.load()
	m.download()
}

func (m *RepoBuild) load() {
	fp := filepath.Join(m.WD, "repo.yml")
	err := marshal.FromFile(&m.Serialized, fp)
	m.ErrHandler.FatalF("Can not load file \"%s\": %v", fp, err)

	err = m.Validate.Struct(&m.Serialized)
	m.ErrHandler.FatalF("Can not load file \"%s\", invalid fields: %v", fp, err)
}

func (m *RepoBuild) download() {
	destDir := filepath.Join(m.WD, "templates")
	err := os.MkdirAll(destDir, 0755)
	errs.FatalF("Can create directory \"%s\": %v", destDir, err)

	for _, url := range m.Serialized.TemplateURLs {
		// Validate url
		err := m.Validate.Var(url, "url")
		if m.ErrHandler.LogF("Invalid url \"%s\": %v", url, err) {
			continue
		}

		// Get URL
		srcDir, err := gitclient.Get(url, &m.Plumb)
		if m.ErrHandler.LogF("Can not download \"%s\": %v", url, err) {
			continue
		}

		// Load .kick.yml
		var templateMain serialize.TemplateMain
		srcTemplate := filepath.Join(srcDir, ".kick.yml")
		err = marshal.FromFile(&templateMain, srcTemplate)
		if m.ErrHandler.LogF("Can not load file \"%s\": %v", srcTemplate, err) {
			continue
		}

		// Validate .kick.yml
		err = m.Validate.Struct(&templateMain)
		if err != nil {
			var invalid []string
			for _, err := range err.(validator.ValidationErrors) {
				invalid = append(invalid, err.StructField())
			}
			m.Log.Printf("Can not load %s invalid fields: ", strings.Join(invalid, `,`))
			continue
		}

		// Copy object to "templates/*.yml" yaml file
		var templateElement serialize.RepoTemplateFile
		err = copier.Copy(&templateElement, &templateMain)
		if m.ErrHandler.LogF("Can not copy objects: %v", err) {
			continue
		}
		// Add URL
		templateElement.URL = url

		// Write "templates/*.yml" yaml file
		destRepoYAML := filepath.Join(destDir, templateElement.Name+".yml")
		err = marshal.ToFile(&templateElement, destRepoYAML)
		if m.ErrHandler.LogF("Can not save file \"%s\": %v", destRepoYAML, err) {
			continue
		}
	}
}
