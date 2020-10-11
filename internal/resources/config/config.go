package config

import (
	"strings"

	"github.com/crosseyed/prjstart/internal/utils/marshal"
)

var userconfig string = ".prjstart.yml"

// File configuration as loaded from the configuration file
type File struct {
	Path       string     `yaml:"-"`                 // Path to configuration file
	MasterURLs []string   `yaml:"masters,omitempty"` // URLs to master git repositories
	Templates  []Template `yaml:"templates,flow"`    // Template definitions
}

// Options options to New
type Options struct {
	Path string // Path to configuration file
}

// New Config constructor
func New(opts Options) *File {
	if opts.Path == "" {
		panic("opts.ConfigFile can not be an empty string")
	}
	c := &File{
		Path: opts.Path,
	}
	return c
}

// SortByName sort template alphabetically by name
type SortByName []Template

func (a SortByName) Len() int           { return len(a) }
func (a SortByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByName) Less(i, j int) bool { return strings.Compare(a[j].Name, a[i].Name) > 0 }

// Template template configuration in main configuration file
type Template struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Desc string `yaml:"desc"`
}

// Load loads configuration file from disk
func (f *File) Load() {
	marshal.UnmarshalFile(f, f.Path)
}

// Save saves configuration file to disk
func (f *File) Save() {
	marshal.MarshalFile(f, f.Path)
}
