package repocmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/kick-project/kick/internal/subcmds/repocmd"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", repocmd.UsageDoc)
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
    - http://127.0.0.1:8080/tmpl.git
    - http://127.0.0.1:8080/tmpl1.git
    - http://127.0.0.1:8080/tmpl2.git
    - http://127.0.0.1:8080/tmpl3.git
    - http://127.0.0.1:8080/tmpl4.git
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
	inject := di.New(&di.Options{Home: home})
	repocmd.Repo([]string{"repo", "build"}, inject)

	type stru struct {
		Name     string   `yaml:"name"`
		Desc     string   `yaml:"description"`
		URL      string   `yaml:"url"`
		Versions []string `yaml:"versions"`
	}

	for _, id := range []string{"tmpl", "tmpl1", "tmpl2", "tmpl3", "tmpl4"} {
		d := &stru{}
		y := filepath.Join(dirPath, "templates", fmt.Sprintf(`%s.yml`, id))
		err := marshal.FromFile(d, y)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, d.Name, id)
		assert.Equal(t, d.URL, fmt.Sprintf(`http://127.0.0.1:8080/%s.git`, id))
		assert.Equal(t, d.Desc, fmt.Sprintf(`%s template`, id))
		assert.Greater(t, len(d.Versions), 0)
	}
}

func TestRepocmd_List(t *testing.T) {
	home := filepath.Join(testtools.TempDir(), "home")
	inject := di.New(&di.Options{Home: home})
	stdout := bytes.NewBufferString(``)
	inject.Stdout = stdout
	repocmd.Repo([]string{"repo", "list"}, inject)

	mustMatch := regexp.MustCompile(`\| repo1 \| http://127.0.0.1:8080/repo1.git \|`)
	mustNotMatch := regexp.MustCompile(`local`)
	assert.Regexp(t, mustMatch, stdout.String())
	assert.NotRegexp(t, mustNotMatch, stdout.String())
}

func TestRepocmd_Info(t *testing.T) {
	home := filepath.Join(testtools.TempDir(), "home")
	inject := di.New(&di.Options{Home: home})
	stdout := bytes.NewBufferString(``)
	inject.Stdout = stdout
	repocmd.Repo([]string{"repo", "info", "repo1"}, inject)

	mustMatch1 := regexp.MustCompile(`name: repo1`)
	mustMatch2 := regexp.MustCompile(`url: http://127.0.0.1:8080/repo1.git`)
	assert.Regexp(t, mustMatch1, stdout)
	assert.Regexp(t, mustMatch2, stdout)
}
