package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/crosseyed/prjstart/internal/utils/dfaults"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/crosseyed/prjstart/internal/utils/marshal"
	"gopkg.in/yaml.v2"
)

var userconfig string = ".prjstart.yml"
var userconfigdir string = ".prjstart"

// File configuration as loaded from the configuration file
type File struct {
	Path         string         // Path to configuration file
	Home         string         // Home directory
	MasterURLs   []string       `yaml:"masters"`   // URLs to master git repos
	TemplateURLs []TemplateStub `yaml:"templates"` // URLs to template git repos
}

// Options options to New
type Options struct {
	Home string // Home directory
	Path string // Path to configuration file
}

// New Config constructor
func New(opts Options) *File {
	if opts.Path == "" {
		panic("opts.ConfigFile can not be an empty string")
	}
	c := &File{
		Home: opts.Home,
		Path: opts.Path,
	}
	return c
}

// SortByName sort template alphabetically by name
type SortByName []TemplateStub

func (a SortByName) Len() int           { return len(a) }
func (a SortByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByName) Less(i, j int) bool { return strings.Compare(a[j].Name, a[i].Name) > 0 }

// TemplateStub template configuration in main configuration file
type TemplateStub struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Desc string `yaml:"desc"`
}

// Load loads configuration file from disk
func (f *File) Load() {
	marshal.UnmarshalFile(f, f.Path)
}

// Load loads configuration from disk
func Load(homedir, prjstart string) *File {
	homedir = dfaults.String(os.Getenv("HOME"), homedir)
	prjstart = dfaults.String(userconfig, prjstart)
	conffile := filepath.Join(homedir, prjstart)
	if _, err := os.Stat(conffile); os.IsNotExist(err) {
		return nil
	}

	f, err := ioutil.ReadFile(conffile)
	errutils.Efatalf("Can not read file %s: %v", conffile, err)
	conf := File{
		Home: homedir,
	}

	err = yaml.Unmarshal([]byte(f), &conf)
	errutils.Efatalf("Can not unmarshal file %s: %v", conffile, err)
	return &conf
}
