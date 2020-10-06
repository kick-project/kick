package metadata

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	fp "path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/services/initialize"
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"syreclabs.com/go/faker"
)

func TestBuild(t *testing.T) {
	initIt()
	conf := &config.File{
		Home: fp.Join(utils.TempDir(), "home"),
		MasterURLs: []string{
			"http://127.0.0.1:5000/master1.git",
		},
	}
	dbfile := fp.Clean(fmt.Sprintf("%s/home/.prjstart/metadata/metadata.db", utils.TempDir()))
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on", dbfile)
	dbconn, err := sql.Open("sqlite3", dsn)
	errutils.Efatalf("%w", err)
	m := New(Options{
		ConfigFile:   conf,
		MetadataPath: fp.Clean(fmt.Sprintf("%s/home/.prjstart/metadata", utils.TempDir())),
		DB:           dbconn,
	})
	m.Build()
}

func initIt() {
	tmpdir := utils.TempDir()
	confpath := fp.Clean(fmt.Sprintf("%s/home/.prjstart.yml", tmpdir))
	templatedir := fp.Clean(fmt.Sprintf("%s/home/.prjstart/project", tmpdir))
	metadatadir := fp.Clean(fmt.Sprintf("%s/home/.prjstart/metadata", tmpdir))
	metadatadb := fp.Clean(fmt.Sprintf("%s/home/.prjstart/metadata/metadata.db", tmpdir))
	i := initialize.New(initialize.Options{
		ConfigPath:  confpath,
		TemplateDir: templatedir,
		MetadataDir: metadatadir,
		SQLiteFile:  metadatadb,
		DBDriver:    "sqlite3",
		DSN:         fmt.Sprintf("file:%s?_foreign_key=on", metadatadb),
	})
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
