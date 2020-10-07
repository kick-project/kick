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
	confFile        *config.File
	db              *sql.DB
	dbdriver        string
	dbdsn           string
	sqlitedb        string
	home            string
	pathMetadataDir string
	pathTemplateDir string
	pathUserConf    string
}

// GetSettings get settings using the supplied "home" directory option. Any
// Dependency Injection (DI) configuration created by settings is then
// contextualised by the home environment variable. for instance when home is
// set the paths '{{home}}/prjstart.yml",
// "{{home}}/prjstart/metadata/metadata.db", "{{home}}.prjstart/templates" (etc)
// are then factored in when creating dependency injections.
//
// If inialisation is required then the initialize package can be used. For example
//
//   set := GetSettings("/tmp/tmp_home"); init := initialize.New(set.Initialize())
//
// will create the structures under "/tmp/tmp_home"
//
// If home is an empty string then GetSettings defaults to the $HOME environment
// variable. This is also known as production mode.
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
	if s.confFile != nil {
		return s.confFile
	}
	conf := config.New(config.Options{
		Home: s.home,
		Path: s.pathUserConf,
	})
	conf.Load()
	s.confFile = conf
	return s.confFile
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
	db := s.getDB()
	opts.ConfigPath = s.ConfigFile()
	opts.MetadataDir = s.pathMetadataDir
	opts.DB = db
	return opts
}

//
// Misc
//

func (s *Settings) getDB() *sql.DB {
	if s.db == nil {
		db, err := sql.Open(s.dbdriver, s.dbdsn)
		if err != nil {
			panic(err)
		}
		s.db = db
	}
	return s.db
}
