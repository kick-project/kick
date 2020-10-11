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
	_ "github.com/mattn/go-sqlite3" // Driver for database/sql
)

//
// Settings
//

// Settings provides settings for resources & services.
type Settings struct {
	DBDriver string
	DBDsn    string
	Home     string
	// No Colour output when running commands.
	NoColour bool
	// Project name, normally supplied by the start sub command.
	ProjectName     string
	PathMetadataDir string
	PathTemplateDir string
	PathUserConf    string
	SqliteDB        string
	Stderr          io.Writer
	Stdout          io.Writer
	confFile        *config.File
	db              *sql.DB
}

// GetSettings get settings using the supplied "home" directory option. Any
// Dependency Injection (DI) configuration created by settings is then
// contextualized by the home variable. For instance when home is
// set the paths '{{home}}/prjstart.yml",
// "{{home}}/prjstart/metadata/metadata.db", "{{home}}.prjstart/templates" (etc)
// are then factored in when creating dependency injections.
//
// If initialization is needed for testing then the initialize package can be
// used. For example
//
//   set := GetSettings("/tmp/tmp_home");
//   init := initialize.New(iinitialize.Inject(s));
//   init.Init()
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
		DBDriver:        dbdriver,
		DBDsn:           dbdsn,
		SqliteDB:        sqlitedb,
		Home:            home,
		PathMetadataDir: pathMetadataDir,
		PathTemplateDir: pathTemplateDir,
		PathUserConf:    pathUserConf,
		Stderr:          os.Stderr,
		Stdout:          os.Stdout,
	}
	return s
}

//
// Tools - The tools in this section should only be used in an injector or for
// testing purposes.
//

// ConfigFile load settings from configuration file
func (s *Settings) ConfigFile() *config.File {
	if s.confFile != nil {
		return s.confFile
	}
	conf := config.New(config.Options{
		Path: s.PathUserConf,
	})
	conf.Load()
	s.confFile = conf
	return s.confFile
}

// GetDB return DB object.
func (s *Settings) GetDB() *sql.DB {
	if s.db == nil {
		db, err := sql.Open(s.DBDriver, s.DBDsn)
		if err != nil {
			panic(err)
		}
		s.db = db
	}
	return s.db
}
