package initialize

import (
	"database/sql"
	"os"
	fp "path/filepath"

	"github.com/crosseyed/prjstart/internal/resources/db"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

// Init is responsibile for initializing all disk paths
type Init struct {
	confpath    string
	templatedir string
	metadatadir string
	sqlitefile  string
	driver      string
	dsn         string
}

// Options options to New
type Options struct {
	ConfigPath  string // Path to configuration file
	DBDriver    string // SQL Driver to use
	DSN         string // SQL DSN
	MetadataDir string // Path to metadata directory
	SQLiteFile  string // Path to DB file
	TemplateDir string // Path to template directory
}

// New creates a new *Init object which is responsibile for initializing all directory structures
func New(opts Options) *Init {
	i := &Init{
		confpath:    opts.ConfigPath,
		templatedir: opts.TemplateDir,
		metadatadir: opts.MetadataDir,
		sqlitefile:  opts.SQLiteFile,
		driver:      opts.DBDriver,
		dsn:         opts.DSN,
	}
	return i
}

// Init initialize everything.
func (i *Init) Init() {
	i.InitPaths()
	i.InitMetadata()
}

// InitPaths initialize paths.
func (i *Init) InitPaths() {
	confdir := fp.Dir(i.confpath)
	dbdir := fp.Dir(i.sqlitefile)
	mkDirs([]string{confdir, dbdir, i.templatedir, i.metadatadir})
}

// InitMetadata initialize metadata.
func (i *Init) InitMetadata() {
	dbdir := fp.Dir(i.sqlitefile)
	mkDirs(dbdir)
	dbconn, err := sql.Open(i.driver, i.dsn)
	errutils.Epanicf("can not connect to database: %w", err)
	db.CreateSchema(dbconn)
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
