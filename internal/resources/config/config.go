package config

import (
	"strings"

	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/crosseyed/prjstart/internal/utils/marshal"
)

// File configuration as loaded from the configuration file
type File struct {
	pathUserConf     string     `yaml:"-"` // Path to configuration file
	pathTemplateConf string     `yaml:"-"`
	MasterURLs       []string   `yaml:"masters,omitempty"` // URLs to master git repositories
	Templates        []Template `yaml:"-"`                 // Template definitions
}

// Options options to New
type Options struct {
	PathUserConf     string // Path to config.yml
	PathTemplateConf string // Path to configuration file
}

// New Config constructor
func New(opts Options) *File {
	if opts.PathUserConf == "" {
		panic("opts.PathUserConf can not be an empty string")
	}
	if opts.PathTemplateConf == "" {
		panic("opts.PathTemplateConf can not be an empty string")
	}
	c := &File{
		pathUserConf:     opts.PathUserConf,
		pathTemplateConf: opts.PathTemplateConf,
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
	err := marshal.UnmarshalFile(f, f.pathUserConf)
	errutils.Epanicf("%w", err)
	err = marshal.UnmarshalFile(&f.Templates, f.pathTemplateConf)
	errutils.Epanicf("%w", err)
}

// Save saves configuration file to disk
func (f *File) Save() {
	marshal.MarshalFile(f.Templates, f.pathTemplateConf)
}
