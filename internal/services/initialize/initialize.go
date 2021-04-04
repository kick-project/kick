package initialize

import (
	"database/sql"
	"os"
	fp "path/filepath"

	"github.com/kick-project/kick/internal/fflags"
	"github.com/kick-project/kick/internal/resources/db"
	"github.com/kick-project/kick/internal/utils/errutils"
)

// Initialize is responsible for initializing all disk paths
type Initialize struct {
	ConfigPath         string
	ConfigTemplatePath string
	DBDriver           string
	DSN                string
	HomeDir            string
	MetadataDir        string
	SQLiteFile         string
	TemplateDir        string
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
	if fflags.ORM() {
		i.initMetadataNew()
	} else {
		i.initMetadataOld()
	}
}

func (i *Initialize) initMetadataNew() {
	db.CreateModel(&db.ModelOptions{
		File: i.SQLiteFile,
	})
}

func (i *Initialize) initMetadataOld() {
	dbdir := fp.Dir(i.SQLiteFile)
	mkDirs(dbdir)
	dbconn, err := sql.Open(i.DBDriver, i.DSN)
	errutils.Epanicf("can not connect to database: %w", err)
	db.CreateSchema(dbconn)
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
	// TODO: Marshal with welcome content
	//if _, err := os.Stat(i.confpath); os.IsNotExist(err) {
	//	i.configfile.Save()
	//} else if err != nil {
	//	errutils.Epanicf("can not save configuration file: %w", err)
	//}
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
