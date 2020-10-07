package settings

import (
	"database/sql"
	"fmt"
	"os"
	fp "path/filepath"

	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/services/template/variables"
	"github.com/crosseyed/prjstart/internal/utils/dfaults"
	_ "github.com/mattn/go-sqlite3" // Driver for database/sql
)

//
// Settings
//

// Settings provides settings for resources & services.
type Settings struct {
	DBConfig        *DBConfig
	File            *config.File
	db              *sql.DB
	dbdriver        string
	dbdsn           string
	sqlitedb        string
	home            string
	pathMetadataDir string
	pathTemplateDir string
	pathUserConf    string
}

// GetSettings get settings
func GetSettings(home string) *Settings {
	home = dfaults.String(os.Getenv("HOME"), home)
	dbdriver := "sqlite3"
	sqlitedb := fp.Clean(fmt.Sprintf("%s/.prjstart/metadata/metadata.db", home))
	dbdsn := fmt.Sprintf("file:%s?_foreign_key=on", sqlitedb)
	pathUserConf := fp.Clean(fmt.Sprintf("%s/.prjstart.yml", home))
	pathTemplateDir := fp.Clean(fmt.Sprintf("%s/.prjstart/templates", home))
	pathMetadataDir := fp.Clean(fmt.Sprintf("%s/.prjstart/metadata", home))
	s := &Settings{
		dbdriver:        dbdriver,
		dbdsn:           dbdsn,
		sqlitedb:        sqlitedb,
		home:            home,
		pathMetadataDir: pathMetadataDir,
		pathTemplateDir: pathTemplateDir,
		pathUserConf:    pathUserConf,
	}
	return s
}

// ConfigFile load settings from configuration file
func (s *Settings) ConfigFile() *config.File {
	if s.File != nil {
		return s.File
	}
	conf := config.New(config.Options{
		Home: s.home,
		Path: s.pathUserConf,
	})
	conf.Load()
	s.File = conf
	return s.File
}

// Template creates settings for template.New
func (s *Settings) Template() (opts struct {
	Config      *config.File
	Variables   *variables.Variables
	TemplateDir string
	ModeLineLen uint8
}) {
	configFile := s.ConfigFile()
	vars := variables.NewTmplVars()
	opts.Config = configFile
	opts.TemplateDir = s.pathTemplateDir
	opts.Variables = vars

	return opts
}

// Initialize creates settings for initialize.New
func (s *Settings) Initialize() (opts struct {
	ConfigPath  string // Path to configuration file
	DBDriver    string // SQL Driver to use
	DSN         string // SQL DSN
	MetadataDir string // Path to metadata directory
	SQLiteFile  string // Path to DB file
	TemplateDir string // Path to template directory
}) {
	opts.ConfigPath = s.pathUserConf
	opts.DBDriver = s.dbdriver
	opts.DSN = s.dbdsn
	opts.MetadataDir = s.pathMetadataDir
	opts.SQLiteFile = s.sqlitedb
	opts.TemplateDir = s.pathTemplateDir
	return opts
}

// Metadata creates settings for metadata.New
func (s *Settings) Metadata() (opts struct {
	ConfigPath  *config.File
	MetadataDir string
	DB          *sql.DB
}) {
	if s.db == nil {
		db, err := sql.Open(s.dbdriver, s.dbdsn)
		if err != nil {
			panic(err)
		}
		s.db = db
	}
	opts.ConfigPath = s.ConfigFile()
	opts.MetadataDir = s.pathMetadataDir
	opts.DB = s.db
	return opts
}

//
// Database configuration
//

// DBConfig Database configuration
type DBConfig struct {
	SQLitePath string // Path to db file
	Driver     string
	DSN        string
	db         *sql.DB
}

// DB returns the *sql.DB database object
func (s DBConfig) DB() *sql.DB {
	if s.db != nil {
		return s.db
	}
	db, err := sql.Open(s.Driver, s.DSN)
	if err != nil {
		panic(err)
	}
	s.db = db
	return db
}
