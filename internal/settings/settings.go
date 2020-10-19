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

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/kick-project/kick/internal/env"
	"github.com/kick-project/kick/internal/resources/config"
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
	// See https://pkg.go.dev/github.com/apex/log#InfoLevel for
	// available levels.
	logLevel log.Level
	// No Colour output when running commands.
	NoColour bool
	// Project name, normally supplied by the start sub command.
	ProjectName      string
	PathMetadataDir  string
	PathTemplateConf string
	PathTemplateDir  string
	PathUserConf     string
	SqliteDB         string
	Stdin            io.Reader
	Stderr           io.Writer
	Stdout           io.Writer
}

// GetSettings get settings using the supplied "home" directory option. Any
// Dependency Injection (DI) configuration created by settings is then
// contextualized by the home variable. For instance when home is
// set the paths...
//
//     {{home}}/.kick/config.yml
//     {{home}}/.kick/templates.yml
//     {{home}}/.kick/metadata/metadata.db
//     {{home}}/.kick/templates
//     etc..
//
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
	sqlitedb := fp.Clean(fmt.Sprintf("%s/.kick/metadata/metadata.db", home))
	dbdsn := fmt.Sprintf("file:%s?_foreign_key=on", sqlitedb)
	pathUserConf := fp.Clean(fmt.Sprintf("%s/.kick/config.yml", home))
	pathTemplateConf := fp.Clean(fmt.Sprintf("%s/.kick/templates.yml", home))
	pathTemplateDir := fp.Clean(fmt.Sprintf("%s/.kick/templates", home))
	pathMetadataDir := fp.Clean(fmt.Sprintf("%s/.kick/metadata", home))
	logLvl := log.ErrorLevel
	if env.Debug() {
		logLvl = log.DebugLevel
	}
	s := &Settings{
		DBDriver:         dbdriver,
		DBDsn:            dbdsn,
		SqliteDB:         sqlitedb,
		Home:             home,
		PathMetadataDir:  pathMetadataDir,
		PathTemplateConf: pathTemplateConf,
		PathTemplateDir:  pathTemplateDir,
		PathUserConf:     pathUserConf,
		Stderr:           os.Stderr,
		Stdin:            os.Stdin,
		Stdout:           os.Stdout,
		logLevel:         logLvl,
	}
	return s
}

// LogLevel Sets the log level
func (s *Settings) LogLevel(lvl log.Level) {
	s.logLevel = lvl
}

//
// Tools - The tools in this section should only be used in an injector or for
// testing purposes.
//

// ConfigFile load settings from configuration file
func (s *Settings) ConfigFile() *config.File {
	conf := &config.File{
		PathUserConf:     s.PathUserConf,
		PathTemplateConf: s.PathTemplateConf,
	}
	conf.Load()
	return conf
}

// GetDB return DB object.
func (s *Settings) GetDB() *sql.DB {
	db, err := sql.Open(s.DBDriver, s.DBDsn)
	if err != nil {
		panic(err)
	}
	return db
}

// GetLogger inject logger object.
func (s *Settings) GetLogger() *log.Logger {
	logger := &log.Logger{
		Handler: text.New(s.Stderr),
		Level:   s.logLevel,
	}
	return logger
}
