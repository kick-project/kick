// Package settings is a package that implements Dependency Injection through
// methods that create the options needed for structs to be created. See
// GetSettings for more information.
//
// This package uses anonymous structs as the injection options, this is to
// avoid issues with import loops when importing settings into *_test.go files.
package settings

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	fp "path/filepath"

	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/services/template/variables"
	_ "github.com/mattn/go-sqlite3" // Driver for database/sql
)

//
// Settings
//

// Settings provides settings for resources & services.
type Settings struct {
	NoColour        bool
	confFile        *config.File
	db              *sql.DB
	dbdriver        string
	dbdsn           string
	Home            string
	PathMetadataDir string
	PathTemplateDir string
	PathUserConf    string
	sqlitedb        string
	Stderr          io.Writer
	Stdout          io.Writer
}

// GetSettings get settings using the supplied "home" directory option. Any
// Dependency Injection (DI) configuration created by settings is then
// contextualized by the home variable. For instance when home is
// set the paths '{{home}}/prjstart.yml",
// "{{home}}/prjstart/metadata/metadata.db", "{{home}}.prjstart/templates" (etc)
// are then factored in when creating dependency injections.
//
// If when initialization is needed then the initialize package can be used. For example
//
//   set := GetSettings("/tmp/tmp_home"); init := initialize.New(set.Initialize()); init.Init()
//
// will create the structures under "/tmp/tmp_home"
//
// "home" must be explicitly set or a panic will ensue.
func GetSettings(home string) *Settings {
	if home == "" {
		panic("home is set to an empty string")
	}
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
		Home:            home,
		PathMetadataDir: pathMetadataDir,
		PathTemplateDir: pathTemplateDir,
		PathUserConf:    pathUserConf,
		Stderr:          os.Stderr,
		Stdout:          os.Stdout,
	}
	return s
}

// ConfigFile load settings from configuration file
func (s *Settings) ConfigFile() *config.File {
	if s.confFile != nil {
		return s.confFile
	}
	conf := config.New(config.Options{
		Home: s.Home,
		Path: s.PathUserConf,
	})
	conf.Load()
	s.confFile = conf
	return s.confFile
}

//
// Injectors
//

// Initialize creates settings for initialize.New
func (s *Settings) Initialize() (opts struct {
	ConfigFile  *config.File // Initialized config file
	ConfigPath  string       // Path to configuration file
	DBDriver    string       // SQL Driver to use
	DSN         string       // SQL DSN
	HomeDir     string       // Path to home directory
	MetadataDir string       // Path to metadata directory
	SQLiteFile  string       // Path to DB file
	TemplateDir string       // Path to template directory
}) {
	conf := config.New(config.Options{
		Home: s.Home,
		Path: s.PathUserConf,
	})
	opts.ConfigFile = conf
	opts.ConfigPath = s.PathUserConf
	opts.DBDriver = s.dbdriver
	opts.DSN = s.dbdsn
	opts.HomeDir = s.Home
	opts.MetadataDir = s.PathMetadataDir
	opts.SQLiteFile = s.sqlitedb
	opts.TemplateDir = s.PathTemplateDir
	return opts
}

// Metadata creates settings for metadata.New
func (s *Settings) Metadata() (opts struct {
	ConfigFile  *config.File
	MetadataDir string
	DB          *sql.DB
}) {
	db := s.GetDB()
	opts.ConfigFile = s.ConfigFile()
	opts.MetadataDir = s.PathMetadataDir
	opts.DB = db
	return opts
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
	opts.TemplateDir = s.PathTemplateDir
	opts.Variables = vars

	return opts
}

//
// Misc
//

// GetDB return DB object. This should only be used for testing, all other calls to
// the DB object should be performed through an Injector.
func (s *Settings) GetDB() *sql.DB {
	if s.db == nil {
		db, err := sql.Open(s.dbdriver, s.dbdsn)
		if err != nil {
			panic(err)
		}
		s.db = db
	}
	return s.db
}
