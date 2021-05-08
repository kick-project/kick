package repobuild_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/services/repobuild"
	"github.com/kick-project/kick/internal/utils"
)

func TestRepoBuildRepo(t *testing.T) {
	// Make directory
	dirPath := filepath.Clean(utils.TempDir() + "/TestRepoBuildRepo")
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		t.Errorf("Can not create directory %s: %v", dirPath, err)
		return
	}

	// Make file
	filePath := filepath.Clean(dirPath + "/repo.yml")
	data := []byte(
		`name: repo1
description: repo1 repo
templates:
    - http://127.0.0.1:5000/tmpl.git
    - http://127.0.0.1:5000/tmpl1.git
    - http://127.0.0.1:5000/tmpl2.git
    - http://127.0.0.1:5000/tmpl3.git
    - http://127.0.0.1:5000/tmpl4.git
`)

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		t.Errorf("Can not write file %s: %v", filePath, err)
		return
	}

	m := repobuild.RepoBuild{
		WD: dirPath,
		Plumb: plumbing.Plumbing{
			Base: filepath.Join(utils.TempDir(), "home", ".kick", "metadata"),
		},
	}
	m.Make()
}
