package initialize

import (
	"os"
	fp "path/filepath"

	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/utils/errutils"
)

// Initialize is responsible for initializing all disk paths
type Initialize struct {
	ConfigPath         string `copier:"must"`
	ConfigTemplatePath string `copier:"must"`
	DBDriver           string `copier:"must"`
	DSN                string `copier:"must"`
	HomeDir            string `copier:"must"`
	MetadataDir        string `copier:"must"`
	SQLiteFile         string `copier:"must"`
	TemplateDir        string `copier:"must"`
}

// Init initialize everything.
func (i *Initialize) Init() {
	i.InitPaths()
	i.InitMetadata()
	i.InitConfig()
}

// InitPaths initialize paths.
func (i *Initialize) InitPaths() {
	confdir := fp.Dir(i.ConfigPath)
	dbdir := fp.Dir(i.SQLiteFile)
	mkDirs([]string{confdir, dbdir, i.TemplateDir, i.MetadataDir})
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
		errutils.Elogf("error: %w", err)
		defer f.Close()
		_, err = f.WriteString(`---
`)
		errutils.Epanic(err)
	} else if err != nil {
		errutils.Epanicf("can not save configuration file: %w", err)
	}
	_, err = os.Stat(i.ConfigTemplatePath)
	if os.IsNotExist(err) {
		f, err := os.Create(i.ConfigTemplatePath)
		errutils.Elogf("error: %w", err)
		defer f.Close()
		_, err = f.WriteString(`---
`)
		errutils.Epanic(err)
	} else if err != nil {
		errutils.Epanicf("can not save configuration file: %w", err)
	}
}

func mkDirs(i interface{}) {
	var dirs []string
	switch v := i.(type) {
	case string:
		dirs = []string{v}
	case []string:
		dirs = v
	default:
		panic("unknown type")
	}
	for _, d := range dirs {
		err := os.MkdirAll(d, 0755)
		errutils.Epanicf("can not create %s: %w", d, err)
	}
}
