package metadata

import (
	"fmt"
	"io/ioutil"
	"os"
	fp "path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/services/initialize"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/settings/iinitialize"
	"github.com/crosseyed/prjstart/internal/settings/imetadata"
	"github.com/crosseyed/prjstart/internal/utils"
	"syreclabs.com/go/faker"
)

func TestBuild(t *testing.T) {
	home := fp.Join(utils.TempDir(), "home")
	s := settings.GetSettings(home)
	initIt(s)
	m := New(imetadata.Inject(s))
	m.Build()
}

func initIt(s *settings.Settings) {
	i := initialize.New(iinitialize.Inject(s))
	i.Init()
}

func TestMaster_Load(t *testing.T) {
	path, fname, furl, fdesc := fakeJSON(t)
	defer os.Remove(path)
	notEmpty(t, path)
	m := Master{}
	m.Load(path)
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

func TestMaster_Save(t *testing.T) {
	path, _, _, _ := fakeJSON(t)
	defer os.Remove(path)
	m := Master{}
	m.Load(path)
	tmpfile, err := ioutil.TempFile("", "*.json")
	if err != nil {
		t.Fatal("Error opening tempfile")
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())
	m.Save(tmpfile.Name())
	notEmpty(t, tmpfile.Name())
}

func TestTemplate_Load(t *testing.T) {
	path, fname, furl, fdesc := fakeJSON(t)
	defer os.Remove(path)
	notEmpty(t, path)
	tpl := Template{}
	tpl.Load(path)
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
	tpl := Template{}
	tpl.Load(path)
	tmpfile, err := ioutil.TempFile("", "*.json")
	if err != nil {
		t.Fatal("Error opening tempfile")
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())
	tpl.Save(tmpfile.Name())
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
	tf.WriteString(template)
	tf.Close()
	return tf.Name(), name, url, desc
}
