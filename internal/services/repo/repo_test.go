package repo_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-playground/validator"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/kick-project/kick/internal/services/repo"
)

func TestRepo_Build(t *testing.T) {
	// Make directory
	dirPath := filepath.Join(testtools.TempDir(), "TestRepo_Build")
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

	m := repo.Repo{
		WD: dirPath,
		Plumb: &plumbing.Plumbing{
			Base: filepath.Join(testtools.TempDir(), "home", ".kick", "metadata"),
		},
		Validate: validator.New(),
		ErrHandler: &errs.Handler{
			Ex: &exit.Handler{
				Mode: exit.MPanic,
			},
		},
	}
	m.Build()
}