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

// Master build master repo
type Master struct {
	WD         string            // Working Directory
	Plumb      plumbing.Plumbing // Plumbing object
	Serialized serialize.Master  // Serialized config
}

// Build build repo
func (m *Master) Build() {
	m.load()
	m.download()
}

func (m *Master) load() {
	fp := filepath.Join(m.WD, "master.yml")
	err := marshal.UnmarshalFromFile(&m.Serialized, fp)
	errutils.Efatalf("Can not load file \"%s\": %v", fp, err)
}

func (m *Master) download() {
	destDir := filepath.Join(m.WD, "templates")
	err := os.MkdirAll(destDir, 0755)
	errutils.Efatalf("Can create directory \"%s\": %v", destDir, err)

	for _, url := range m.Serialized.TemplateURLs {
		// Get URL
		srcDir, err := gitclient.Get(url, &m.Plumb)
		if errutils.Elogf("Can not download \"%s\": %v", url, err) {
			continue
		}

		// Load master.yml
		var templateSerialize serialize.Kick
		srcMaster := filepath.Join(srcDir, ".kick.yml")
		err = marshal.UnmarshalFromFile(&templateSerialize, srcMaster)
		if errutils.Elogf("Can not load file \"%s\": %v", srcMaster, err) {
			continue
		}

		// Copy object
		var masterElement serialize.MasterElement
		err = copier.Copy(&masterElement, &templateSerialize)
		if errutils.Elogf("Can not copy objects: %v", err) {
			continue
		}

		// Save element.yml
		masterElement.URL = url
		destMasterYAML := filepath.Join(destDir, masterElement.Name+".yml")
		err = marshal.Marshal2File(&masterElement, destMasterYAML)
		if errutils.Elogf("Can not save file \"%s\": %v", destMasterYAML, err) {
			continue
		}
	}
}
