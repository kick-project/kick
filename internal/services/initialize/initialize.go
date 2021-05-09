package initialize

import (
	"os"
	fp "path/filepath"

	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/model"
)

// Initialize is responsible for initializing all disk paths
type Initialize struct {
	ConfigPath         string `validate:"required"`
	ConfigTemplatePath string `validate:"required"`
	HomeDir            string `validate:"required"`
	MetadataDir        string `validate:"required"`
	SQLiteFile         string `validate:"required"`
	TemplateDir        string `validate:"required"`
}

// Init initialize everything.
func (i *Initialize) Init() {
	i.InitPaths()
	i.InitMetadata()
	i.InitConfig()
}

// InitPaths initialize paths.
func (i *Initialize) InitPaths() {
	for _, cur := range []string{fp.Dir(i.ConfigPath), fp.Dir(i.SQLiteFile), i.TemplateDir, i.MetadataDir} {
		if _, err := os.Stat(cur); os.IsNotExist(err) {
			err := os.MkdirAll(cur, 0755)
			errs.PanicF("can not create %s: %w", cur, err)
		}
	}
}

// InitMetadata initialize metadata.
func (i *Initialize) InitMetadata() {
	// Creating an ORM based model
	model.CreateModel(&model.Options{
		File: i.SQLiteFile,
	})
}

// InitConfig initialize configuration file.
func (i *Initialize) InitConfig() {
	_, err := os.Stat(i.ConfigPath)
	if os.IsNotExist(err) {
		f, err := os.Create(i.ConfigPath)
		errs.LogF("error: %w", err)
		defer f.Close()
		_, err = f.WriteString(`---
`)
		errs.Panic(err)
	} else if err != nil {
		errs.PanicF("can not save configuration file: %w", err)
	}
	_, err = os.Stat(i.ConfigTemplatePath)
	if os.IsNotExist(err) {
		f, err := os.Create(i.ConfigTemplatePath)
		errs.LogF("error: %w", err)
		defer f.Close()
		_, err = f.WriteString(`---
`)
		errs.Panic(err)
	} else if err != nil {
		errs.PanicF("can not save configuration file: %w", err)
	}
}
