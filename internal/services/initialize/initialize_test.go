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
	err := os.MkdirAll(home, 0755)
	errs.Panic(err)

	inject := di.Setup(home)

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
	init.CreateRepo(repo, repo)

	assert.DirExists(t, repoPath)
	assert.FileExists(t, repoYAML)
}

// TestInitialize_Repo tests repo initialization
func TestInitialize_Repo_NoPath(t *testing.T) {
	home := filepath.Join(testtools.TempDir(), "TestInitialize_Repo_NoPath")
	err := os.MkdirAll(home, 0755)
	errs.Panic(err)

	inject := di.Setup(home)

	repo := `myrepo`
	repoPath := filepath.Join(home, repo)
	err = os.MkdirAll(repoPath, 0755)
	errs.Panic(err)
	repoYAML := filepath.Join(repoPath, "repo.yml")

	wd, err := os.Getwd()
	errs.Panic(err)

	err = os.Chdir(repoPath)
	errs.Panic(err)
	defer func() {
		err = os.Chdir(wd)
		errs.Panic(err)
	}()
	init := inject.MakeInit()
	init.CreateRepo(repo, "")

	assert.DirExists(t, repoPath)
	assert.FileExists(t, repoYAML)
}

// TestInitialize_Template test template creation
func TestInitialize_Template(t *testing.T) {
	home := filepath.Join(testtools.TempDir(), "TestInitialize_Template")
	err := os.MkdirAll(home, 0755)
	errs.Panic(err)

	inject := di.Setup(home)

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
	init.CreateTemplate(tmpl, tmpl)

	assert.DirExists(t, tmplPath)
	assert.FileExists(t, tmplYAML)
}

// TestInitialize_Template_NoDir test template initialization
func TestInitialize_Template_NoDir(t *testing.T) {
	home := filepath.Join(testtools.TempDir(), "TestInitialize_Template_NoDir")
	err := os.MkdirAll(home, 0755)
	errs.Panic(err)

	inject := di.Setup(home)

	tmpl := `mytemplate`
	tmplPath := filepath.Join(home, tmpl)
	err = os.MkdirAll(tmplPath, 0755)
	errs.Panic(err)
	tmplYAML := filepath.Join(tmplPath, ".kick.yml")

	wd, err := os.Getwd()
	errs.Panic(err)

	err = os.Chdir(tmplPath)
	errs.Panic(err)
	defer func() {
		err = os.Chdir(wd)
		errs.Panic(err)
	}()
	init := inject.MakeInit()
	init.CreateTemplate(tmpl, "")

	assert.DirExists(t, tmplPath)
	assert.FileExists(t, tmplYAML)
}
