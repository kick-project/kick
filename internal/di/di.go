// Package di is a package that implements Dependency Injection through
// methods that create the options needed for structs to be created. See
// GetDI for more information.
//
// This package uses anonymous structs as the injection options, this is to
// avoid issues with import loops when importing di into *_test.go files.
package di

import (
	"fmt"
	"io"
	"log"
	"os"
	fp "path/filepath"

	apexlog "github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/kick-project/kick/internal/env"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/utils/errutils"
	_ "github.com/mattn/go-sqlite3" // Driver for database/sql
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

//
// DI
//

// DI provides Dependency Injection Container for resources & services. Injectors are in secperate pacages below this package
// to get around dependency loops when injecting into unit tests.
type DI struct {
	Home string
	// See https://pkg.go.dev/github.com/apex/log#InfoLevel for
	// available levels.
	logLevel apexlog.Level

	// Standard logging
	StdLogFlags  int    // Go standard logging flags
	StdLogPrefix string // Go standard logging prefix

	ExitMode int // Exit mode either exit.None or exit.Panic. See package exit for context

	// No Colour output when running commands.
	NoColour bool
	// Project name, normally supplied by the start sub command.
	ProjectName      string
	PathMetadataDir  string
	PathTemplateConf string
	PathRepoDir      string
	PathTemplateDir  string
	PathUserConf     string
	SqliteDB         string
	ModelDB          string
	Stdin            io.Reader
	Stderr           io.Writer
	Stdout           io.Writer
}

// Setup get di using the supplied "home" directory option. Any
// Dependency Injection (DI) configuration created by di is then
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
//   set := Setup("/tmp/tmp_home");
//   init := initialize.New(iinitialize.Inject(s));
//   init.Init()
//
// will create the structures under "/tmp/tmp_home"
//
// "home" must be explicitly set or a panic will ensue.
func Setup(home string) *DI {
	if home == "" {
		panic("home is set to an empty string")
	}
	sqlitedb := fp.Clean(fmt.Sprintf("%s/.kick/metadata/metadata.db", home))
	pathUserConf := fp.Clean(fmt.Sprintf("%s/.kick/config.yml", home))
	pathTemplateConf := fp.Clean(fmt.Sprintf("%s/.kick/templates.yml", home))
	pathRepoDir := fp.Clean(fmt.Sprintf("%s/.kick/repos", home))
	pathTemplateDir := fp.Clean(fmt.Sprintf("%s/.kick/templates", home))
	pathMetadataDir := fp.Clean(fmt.Sprintf("%s/.kick/metadata", home))
	logLvl := apexlog.ErrorLevel
	if env.Debug() {
		logLvl = apexlog.DebugLevel
	}
	s := &DI{
		SqliteDB:         sqlitedb,
		Home:             home,
		PathMetadataDir:  pathMetadataDir,
		PathTemplateConf: pathTemplateConf,
		PathRepoDir:      pathRepoDir,
		PathTemplateDir:  pathTemplateDir,
		PathUserConf:     pathUserConf,
		Stderr:           os.Stderr,
		Stdin:            os.Stdin,
		Stdout:           os.Stdout,
		logLevel:         logLvl,
		StdLogFlags:      log.LstdFlags,
		StdLogPrefix:     "",
		ExitMode:         exit.MNone,
	}
	return s
}

// LogLevel Sets the log level
func (s *DI) LogLevel(lvl apexlog.Level) {
	s.logLevel = lvl
}

//
// Tools - The tools in this section should only be used in an injector or for
// testing purposes.
//

// ConfigFile load di from configuration file
func (s *DI) ConfigFile() *config.File {
	conf := &config.File{
		PathUserConf:     s.PathUserConf,
		PathTemplateConf: s.PathTemplateConf,
	}
	err := conf.Load()
	errutils.Epanic(err)
	return conf
}

// GetORM return ORM object.
func (s *DI) GetORM() *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)
	if _, err = os.Stat(s.SqliteDB); err == nil {
		db, err = gorm.Open(sqlite.Open(s.SqliteDB), &gorm.Config{
			NamingStrategy: &schema.NamingStrategy{
				SingularTable: true,
			},
		})
		errutils.Efatalf("Can not open ORM database %s: %v", s.SqliteDB, err)

	}
	return db
}

// GetLogger inject logger object.
func (s *DI) GetLogger() *apexlog.Logger {
	logger := &apexlog.Logger{
		Handler: text.New(s.Stderr),
		Level:   s.logLevel,
	}
	return logger
}

// GetStdLogger inject Go's standard logging library
func (s *DI) GetStdLogger() *log.Logger {
	logger := log.New(
		s.Stderr,
		"",
		s.StdLogFlags,
	)

	return logger
}

// GetErrorHandler inject error handler
func (s *DI) GetErrorHandler() *errs.Errors {
	handler := errs.Errors{
		Logger: s.GetStdLogger(),
		Ex:     s.GetExitHandler(),
	}
	return &handler
}

// GetExitHandler inject exit handler
func (s *DI) GetExitHandler() *exit.Handler {
	handler := exit.Handler{
		Mode: s.ExitMode,
	}
	return &handler
}
