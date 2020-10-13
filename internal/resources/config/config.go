package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/crosseyed/prjstart/internal/utils/marshal"
)

// File configuration as loaded from the configuration file
type File struct {
	stderr           io.Writer  `yaml:"-"`
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
		stderr:           os.Stderr,
	}
	return c
}

// SortByName sort template alphabetically by name
type SortByName []Template

func (a SortByName) Len() int           { return len(a) }
func (a SortByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByName) Less(i, j int) bool { return strings.Compare(a[j].Handle, a[i].Handle) > 0 }

// Template template configuration in main configuration file
type Template struct {
	Handle   string `yaml:"handle"`
	Template string `yaml:"template"`
	Origin   string `yaml:"origin"`
	URL      string `yaml:"url"`
	Desc     string `yaml:"desc"`
}

// AppendTemplate appends a template to list of templates.
// If stop is non zero, the calling function should exit the program with the
// value contained in stop.
func (f *File) AppendTemplate(t Template) (stop int) {
	// Check if template is installed
	for _, cur := range f.Templates {
		if t.Handle == cur.Handle {
			fmt.Fprintf(os.Stderr, "template handle %s already in use\n", t.Handle)
			return 255
		}
	}

	f.Templates = append(f.Templates, t)
	return 0
}

// Load loads configuration file from disk
func (f *File) Load() {
	err := marshal.UnmarshalFile(f, f.pathUserConf)
	errutils.Epanicf("%w", err)
	err = marshal.UnmarshalFile(&f.Templates, f.pathTemplateConf)
	errutils.Epanicf("%w", err)
}

// SaveTemplates saves template configuration file to disk
func (f *File) SaveTemplates() {
	marshal.MarshalFile(f.Templates, f.pathTemplateConf)
}
