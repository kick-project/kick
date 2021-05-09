package update_test

import (
	"fmt"
	"io/ioutil"
	"os"
	fp "path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/services/update"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"syreclabs.com/go/faker"
)

func TestBuild(t *testing.T) {
	home := fp.Join(utils.TempDir(), "home")
	s := di.Setup(home)
	initIt(s)
	m := s.MakeUpdate()
	err := m.Build()
	if err != nil {
		t.Error(err)
	}
}

func initIt(inject *di.DI) {
	i := inject.MakeInitialize()
	i.Init()
}

func TestRepo_Load(t *testing.T) {
	path, fname, furl, fdesc := fakeJSON(t)
	defer os.Remove(path)
	notEmpty(t, path)
	m := update.Repo{}
	err := m.Load(path)
	if err != nil {
		errutils.Epanic(err)
	}
	if m.Name != fname {
		t.Fail()
	}
	if m.URL != furl {
		t.Fail()
	}
	if m.Description != fdesc {
		t.Fail()
	}
}

func TestRepo_Save(t *testing.T) {
	path, _, _, _ := fakeJSON(t)
	defer os.Remove(path)
	m := update.Repo{}
	err := m.Load(path)
	errutils.Epanic(err)
	tmpfile, err := ioutil.TempFile("", "*.json")
	if err != nil {
		t.Fatal("Error opening tempfile")
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())
	err = m.Save(tmpfile.Name())
	if err != nil {
		t.Error(err)
	}

	notEmpty(t, tmpfile.Name())
}

func TestTemplate_Load(t *testing.T) {
	path, fname, furl, fdesc := fakeJSON(t)
	defer os.Remove(path)
	notEmpty(t, path)
	tpl := update.Template{}
	err := tpl.Load(path)
	if err != nil {
		t.Error(err)
	}
	if tpl.Name != fname {
		t.Fail()
	}
	if tpl.URL != furl {
		t.Fail()
	}
	if tpl.Description != fdesc {
		t.Fail()
	}
}

func TestTemplate_Save(t *testing.T) {
	path, _, _, _ := fakeJSON(t)
	defer os.Remove(path)
	tpl := update.Template{}
	err := tpl.Load(path)
	if err != nil {
		t.Error(err)
	}
	tmpfile, err := ioutil.TempFile("", "*.json")
	if err != nil {
		t.Fatal("Error opening tempfile")
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())
	err = tpl.Save(tmpfile.Name())
	if err != nil {
		t.Error(err)
	}
	notEmpty(t, tmpfile.Name())
}

func notEmpty(t *testing.T, path string) {
	finfo, err := os.Stat(path)
	if err != nil {
		t.Fatalf("File stat error %s: %v", path, err)
	}

	if finfo.Size() == 0 {
		t.Fail()
	}
}

func fakeJSON(t *testing.T) (path, name, url, desc string) {
	name = faker.App().Name()
	url = faker.Internet().Url()
	desc = faker.Internet().Slug()
	template := fmt.Sprintf(`{"name": "%s", "URL": "%s", "description": "%s"}`, name, url, desc)

	tf, err := ioutil.TempFile("", "*.json")
	if err != nil {
		t.Fatalf("ERROR: Can not open temporary file: %v", err)
	}
	_, err = tf.WriteString(template)
	if err != nil {
		t.Error(err)
	}
	tf.Close()
	return tf.Name(), name, url, desc
}
