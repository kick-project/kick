package initialize_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/stretchr/testify/assert"
)

// TestInitialize_Repo tests repo creation
func TestInitialize_Repo(t *testing.T) {
	home := filepath.Join(testtools.TempDir(), "TestInitialize_Repo")
	inject := di.Setup(home)

	err := os.MkdirAll(home, 0755)
	errs.Panic(err)

	repo := `myrepo`
	repoPath := filepath.Join(home, repo)
	repoYAML := filepath.Join(repoPath, "repo.yml")

	wd, err := os.Getwd()
	errs.Panic(err)

	err = os.Chdir(home)
	errs.Panic(err)
	defer func() {
		err = os.Chdir(wd)
		errs.Panic(err)
	}()
	init := inject.MakeInit()
	init.CreateRepo(repo)

	assert.DirExists(t, repoPath)
	assert.FileExists(t, repoYAML)
}

// TestInitialize_Template test template creation
func TestInitialize_Template(t *testing.T) {
	home := filepath.Join(testtools.TempDir(), "TestInitialize_Template")
	inject := di.Setup(home)

	err := os.MkdirAll(home, 0755)
	errs.Panic(err)

	tmpl := `mytemplate`
	tmplPath := filepath.Join(home, tmpl)
	tmplYAML := filepath.Join(tmplPath, ".kick.yml")

	wd, err := os.Getwd()
	errs.Panic(err)

	err = os.Chdir(home)
	errs.Panic(err)
	defer func() {
		err = os.Chdir(wd)
		errs.Panic(err)
	}()
	init := inject.MakeInit()
	init.CreateTemplate(tmpl)

	assert.DirExists(t, tmplPath)
	assert.FileExists(t, tmplYAML)
}
