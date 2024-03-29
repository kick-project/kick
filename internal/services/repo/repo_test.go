package repo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/testtools"
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
    - http://127.0.0.1:8080/tmpl.git
    - http://127.0.0.1:8080/tmpl1.git
    - http://127.0.0.1:8080/tmpl2.git
    - http://127.0.0.1:8080/tmpl3.git
    - http://127.0.0.1:8080/tmpl4.git
`)

	err = os.WriteFile(filePath, data, 0644)
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
	inject := di.New(
		&di.Options{
			Home: home,
		},
	)
	m := inject.MakeRepo()
	m.Build()
}
