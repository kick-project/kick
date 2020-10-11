package initialize

import (
	"database/sql"
	"os"
	fp "path/filepath"

	"github.com/crosseyed/prjstart/internal/resources/db"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

// Initialize is responsibile for initializing all disk paths
type Initialize struct {
	confpath         string
	conftemplatepath string
	driver           string
	dsn              string
	homedir          string
	metadatadir      string
	sqlitefile       string
	templatedir      string
}

// Options options to New
type Options struct {
	ConfigPath         string // Path to configuration file
	ConfigTemplatePath string // Path to template configuration
	DBDriver           string // SQL Driver to use
	DSN                string // SQL DSN
	HomeDir            string // Path to home directory
	MetadataDir        string // Path to metadata directory
	SQLiteFile         string // Path to DB file
	TemplateDir        string // Path to template directory
}

// New creates a new *Init object which is responsibile for initializing all directory structures
func New(opts Options) *Initialize {
	init := &Initialize{
		confpath:         opts.ConfigPath,
		conftemplatepath: opts.ConfigTemplatePath,
		templatedir:      opts.TemplateDir,
		homedir:          opts.HomeDir,
		metadatadir:      opts.MetadataDir,
		sqlitefile:       opts.SQLiteFile,
		driver:           opts.DBDriver,
		dsn:              opts.DSN,
	}
	return init
}

// Init initialize everything.
func (i *Initialize) Init() {
	i.InitPaths()
	i.InitMetadata()
	i.InitConfig()
}

// InitPaths initialize paths.
func (i *Initialize) InitPaths() {
	confdir := fp.Dir(i.confpath)
	dbdir := fp.Dir(i.sqlitefile)
	mkDirs([]string{confdir, dbdir, i.templatedir, i.metadatadir})
}

// InitMetadata initialize metadata.
func (i *Initialize) InitMetadata() {
	dbdir := fp.Dir(i.sqlitefile)
	mkDirs(dbdir)
	dbconn, err := sql.Open(i.driver, i.dsn)
	errutils.Epanicf("can not connect to database: %w", err)
	db.CreateSchema(dbconn)
}

// InitConfig initialize configuration file.
func (i *Initialize) InitConfig() {
	_, err := os.Stat(i.confpath)
	if os.IsNotExist(err) {
		f, err := os.Create(i.confpath)
		errutils.Elogf("error: %w", err)
		defer f.Close()
		_, err = f.WriteString(`---
`)
	} else if err != nil {
		errutils.Epanicf("can not save configuration file: %w", err)
	}
	_, err = os.Stat(i.conftemplatepath)
	if os.IsNotExist(err) {
		f, err := os.Create(i.conftemplatepath)
		errutils.Elogf("error: %w", err)
		defer f.Close()
		_, err = f.WriteString(`---
`)
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
