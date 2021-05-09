package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kick-project/kick/internal/resources/marshal"
)

// File configuration as loaded from the configuration file
type File struct {
	PathTemplateConf string     `yaml:"-"`
	PathUserConf     string     `yaml:"-"` // Path to configuration file
	Stderr           io.Writer  `yaml:"-"`
	RepoURLs         []string   `yaml:"repos,omitempty"` // URLs to repo git repositories
	Templates        []Template `yaml:"-"`               // Template definitions
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
func (f *File) AppendTemplate(t Template) (err error) {
	// Check if template is installed
	for _, cur := range f.Templates {
		if t.Handle == cur.Handle {
			return fmt.Errorf("template handle %s already in use", t.Handle)
		}
	}

	f.Templates = append(f.Templates, t)
	return nil
}

// Load loads configuration file from disk
func (f *File) Load() error {
	pathUserConf := f.PathUserConf
	pathTemplateConf := f.PathTemplateConf
	stderr := f.Stderr

	// Workaround for yaml.v2 clobbering fields with yaml:"-" set.
	// This bug is hard to reproduce as it seems to be intermittent.
	defer func() {
		f.PathUserConf = pathUserConf
		f.PathTemplateConf = pathTemplateConf
		f.Stderr = stderr
	}()

	if _, err := os.Stat(pathUserConf); err == nil {
		err := marshal.FromFile(f, pathUserConf)
		if err != nil {
			return fmt.Errorf("can not load file %s: %w", pathUserConf, err)
		}
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("can not open file %s: %v", pathUserConf, err)
	}

	if _, err := os.Stat(pathTemplateConf); err == nil {
		err = marshal.FromFile(&f.Templates, pathTemplateConf)
		if err != nil {
			return fmt.Errorf("can not load file %s: %w", pathTemplateConf, err)
		}
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("can not open file %s: %w", pathTemplateConf, err)
	}
	return nil
}

// SaveTemplates saves template configuration file to disk
func (f *File) SaveTemplates() error {
	err := marshal.ToFile(f.Templates, f.PathTemplateConf)
	if err != nil {
		return fmt.Errorf("can not save file %s: %w", f.PathTemplateConf, err)
	}
	return nil
}
