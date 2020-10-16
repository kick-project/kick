package config

import (
	"fmt"
	"io"
	"strings"

	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/crosseyed/prjstart/internal/utils/marshal"
)

// File configuration as loaded from the configuration file
type File struct {
	PathTemplateConf string     `yaml:"-"`
	PathUserConf     string     `yaml:"-"` // Path to configuration file
	Stderr           io.Writer  `yaml:"-"`
	MasterURLs       []string   `yaml:"masters,omitempty"` // URLs to master git repositories
	Templates        []Template `yaml:"-"`                 // Template definitions
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
			fmt.Fprintf(f.Stderr, "template handle %s already in use\n", t.Handle)
			return 255
		}
	}

	f.Templates = append(f.Templates, t)
	return 0
}

// Load loads configuration file from disk
func (f *File) Load() {
	err := marshal.UnmarshalFile(f, f.PathUserConf)
	errutils.Efatalf("Can not load file %s: %w\n", f.PathUserConf, err)
	err = marshal.UnmarshalFile(&f.Templates, f.PathTemplateConf)
	errutils.Efatalf("Can not load file %s: %w\n", f.PathTemplateConf, err)
}

// SaveTemplates saves template configuration file to disk
func (f *File) SaveTemplates() {
	err := marshal.MarshalFile(f.Templates, f.PathTemplateConf)
	errutils.Efatalf("Can not save file %s: %w\n", f.PathTemplateConf, err)
}
