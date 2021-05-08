package repobuild

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/gitclient"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/resources/serialize"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/marshal"
)

// RepoBuild build a repository repo
type RepoBuild struct {
	WD         string             // Working Directory
	Plumb      plumbing.Plumbing  // Plumbing object
	Serialized serialize.RepoMain // Serialized config
}

// Make build repo
func (m *RepoBuild) Make() {
	m.load()
	m.download()
}

func (m *RepoBuild) load() {
	fp := filepath.Join(m.WD, "repo.yml")
	err := marshal.FromFile(&m.Serialized, fp)
	errutils.Efatalf("Can not load file \"%s\": %v", fp, err)
}

func (m *RepoBuild) download() {
	destDir := filepath.Join(m.WD, "templates")
	err := os.MkdirAll(destDir, 0755)
	errutils.Efatalf("Can create directory \"%s\": %v", destDir, err)

	for _, url := range m.Serialized.TemplateURLs {
		// Get URL
		srcDir, err := gitclient.Get(url, &m.Plumb)
		if errutils.Elogf("Can not download \"%s\": %v", url, err) {
			continue
		}

		// Load repo.yml
		var templateMain serialize.TemplateMain
		srcTemplate := filepath.Join(srcDir, ".kick.yml")
		err = marshal.FromFile(&templateMain, srcTemplate)
		if errutils.Elogf("Can not load file \"%s\": %v", srcTemplate, err) {
			continue
		}

		// Copy object
		var templateElement serialize.Template
		err = copier.Copy(&templateElement, &templateMain)
		if errutils.Elogf("Can not copy objects: %v", err) {
			continue
		}

		// Save element.yml
		templateElement.URL = url
		destRepoYAML := filepath.Join(destDir, templateElement.Name+".yml")
		err = marshal.ToFile(&templateElement, destRepoYAML)
		if errutils.Elogf("Can not save file \"%s\": %v", destRepoYAML, err) {
			continue
		}
	}
}
