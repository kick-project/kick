package start_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/kick-project/kick/internal/services/start"
	"github.com/stretchr/testify/assert"
)

func TestStart_List_Short(t *testing.T) {
	s, _, stdout := make()
	s.List(false)

	out := stdout.String()
	assert.Contains(t, out, "handle1")
	assert.Contains(t, out, "handle2")
}

func TestStart_List_Long(t *testing.T) {
	s, _, stdout := make()
	s.List(true)
	out := stdout.String()
	assert.Regexp(t, `\|\s+HANDLE\s+\|\s+TEMPLATE\s+\|\s+DESCRIPTION\s+\|\s+LOCATION\s+\|`, out)
	assert.Regexp(t, `\|\s+handle1\s+\|\s+template1/origin1\s+\|\s+-\s+\|\s+http://\S+`, out)
	assert.Regexp(t, `\|\s+handle2\s+\|\s+template2/origin1\s+\|\s+-\s+\|\s+http://\S+`, out)
	assert.Regexp(t, `\|\s+handle3\s+\|\s+template3\s+\|\s+-\s+\|\s+http://\S+`, out)
	assert.Regexp(t, `\|\s+handle4\s+\|\s+-\s+\|\s+-\s+\|\s+http://\S+`, out)
}

func TestStart_Start(t *testing.T) {
	tmpdir := testtools.TempDir()
	path, _ := ioutil.TempDir(tmpdir, "start-")
	project := "tmpl"
	tmpl := "tmpl"

	err := os.Remove(path)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fail()
		}
	}()
	s, _, _ := make()
	s.Start(project, tmpl, path)

	type interpolated struct {
		Project string `yaml:"project"`
		Home    string `yaml:"home"`
	}
	d := &interpolated{}
	y := filepath.Join(path, "template-interpolate.yml")
	err = marshal.FromFile(d, y)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, d.Project, project)
	assert.NotEmpty(t, os.Getenv("HOME"))
	assert.Equal(t, d.Home, os.Getenv("HOME"))
}

func make() (s *start.Start, stderr *bytes.Buffer, stdout *bytes.Buffer) {
	home, _ := filepath.Abs(filepath.Join(testtools.TempDir(), "home"))
	stderr, stdout, conf := getOptions()
	inject := di.New(&di.Options{
		Home: home,
	})
	setup := inject.MakeSetup()
	setup.Init()
	o := start.Options{
		Conf:      conf,
		Check:     inject.MakeCheck(),
		CheckVars: inject.MakeCheckVars(),
		Exit:      &exit.Handler{Mode: exit.MPanic},
		Stderr:    stderr,
		Stdout:    stdout,
		Sync:      inject.MakeSync(),
		Template:  inject.MakeTemplate(),
	}
	s = start.New(o)
	return s, stderr, stdout
}

func getOptions() (stderr, stdout *bytes.Buffer, conf *config.File) {
	stderr = &bytes.Buffer{}
	stdout = &bytes.Buffer{}
	templates := []config.Template{
		{
			Handle:   "handle1",
			Template: "template1",
			Origin:   "origin1",
			URL:      "http://template.io/template1.git",
		},
		{
			Handle:   "handle2",
			Template: "template2",
			Origin:   "origin1",
			URL:      "http://template.io/template2.git",
		},
		{
			Handle:   "handle3",
			Template: "template3",
			URL:      "http://template.io/template3.git",
		},
		{
			Handle: "handle4",
			URL:    "http://template.io/template4.git",
		},
	}
	conf = &config.File{
		Stderr:    stderr,
		Templates: templates,
	}
	return
}
