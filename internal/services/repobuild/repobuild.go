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

// Repo build repo repo
type Repo struct {
	WD         string            // Working Directory
	Plumb      plumbing.Plumbing // Plumbing object
	Serialized serialize.Repo  // Serialized config
}

// Build build repo
func (m *Repo) Build() {
	m.load()
	m.download()
}

func (m *Repo) load() {
	fp := filepath.Join(m.WD, "repo.yml")
	err := marshal.UnmarshalFromFile(&m.Serialized, fp)
	errutils.Efatalf("Can not load file \"%s\": %v", fp, err)
}

func (m *Repo) download() {
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
		var templateSerialize serialize.Kick
		srcRepo := filepath.Join(srcDir, ".kick.yml")
		err = marshal.UnmarshalFromFile(&templateSerialize, srcRepo)
		if errutils.Elogf("Can not load file \"%s\": %v", srcRepo, err) {
			continue
		}

		// Copy object
		var repoElement serialize.RepoElement
		err = copier.Copy(&repoElement, &templateSerialize)
		if errutils.Elogf("Can not copy objects: %v", err) {
			continue
		}

		// Save element.yml
		repoElement.URL = url
		destRepoYAML := filepath.Join(destDir, repoElement.Name+".yml")
		err = marshal.Marshal2File(&repoElement, destRepoYAML)
		if errutils.Elogf("Can not save file \"%s\": %v", destRepoYAML, err) {
			continue
		}
	}
}
