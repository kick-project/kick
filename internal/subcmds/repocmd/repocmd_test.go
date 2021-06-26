package repocmd_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/kick-project/kick/internal/subcmds/repocmd"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", repocmd.UsageDoc)
	assert.NotRegexp(t, "\t", repocmd.UsageDocExtended)
}

func TestRepocmd_Build(t *testing.T) {
	// Make directory
	dirPath := filepath.Join(testtools.TempDir(), "TestRepocmd_Build")
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		t.Errorf("Can not create directory %s: %v", dirPath, err)
		return
	}

	// Make file
	filePath := filepath.Join(dirPath, "repo.yml")
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

	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	_ = os.Chdir(dirPath)
	defer func() { _ = os.Chdir(wd) }()
	home := filepath.Join(testtools.TempDir(), "home")
	inject := di.Setup(home)
	repocmd.Repo([]string{"repo", "build"}, inject)
}

func TestRepocmd_List(t *testing.T) {
	home := filepath.Join(testtools.TempDir(), "home")
	inject := di.Setup(home)
	repocmd.Repo([]string{"repo", "list"}, inject)
}
