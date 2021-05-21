// Package di is a package that implements Dependency Injection through
// methods that create the options needed for structs to be created.
package di

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	fp "path/filepath"
	"runtime"

	"github.com/go-playground/validator"
	"github.com/kick-project/kick/internal/env"
	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/resources/logger"
	"github.com/kick-project/kick/internal/resources/sync"
	"github.com/kick-project/kick/internal/resources/template"
	"github.com/kick-project/kick/internal/resources/template/renderer"
	"github.com/kick-project/kick/internal/resources/template/variables"
	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/services/install"
	"github.com/kick-project/kick/internal/services/list"
	"github.com/kick-project/kick/internal/services/remove"
	"github.com/kick-project/kick/internal/services/repo"
	"github.com/kick-project/kick/internal/services/search"
	"github.com/kick-project/kick/internal/services/setup"
	"github.com/kick-project/kick/internal/services/update"
	_ "github.com/mattn/go-sqlite3" // Driver for database/sql
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

//
// DI
//

// DI provides Dependency Injection Container for resources & services.
type DI struct {
	Home string
	// See https://pkg.go.dev/github.com/apex/log#InfoLevel for
	// available levels.
	logLevel logger.Level

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

	// Cached objects
	cacheConfigFile       *config.File
	cacheORM              *gorm.DB
	cacheErrHandler       *errs.Handler
	cacheExitHandler      *exit.Handler
	cacheCheck            *check.Check
	cacheSetup            *setup.Setup
	cacheList             *list.List
	cacheLogFile          *os.File
	cacheInit             *initialize.Init
	cacheInstall          *install.Install
	cachePlumbingRepo     *plumbing.Plumbing
	cachePlumbingTemplate *plumbing.Plumbing
	cacheRemove           *remove.Remove
	cacheSearch           *search.Search
	cacheSync             *sync.Sync
	cacheTemplate         *template.Template
	cacheUpdate           *update.Update
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
//   init := set.MakeInitialize()
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
	logLvl := logger.ErrorLevel
	if env.Debug() {
		logLvl = logger.DebugLevel
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
func (s *DI) LogLevel(lvl logger.Level) {
	s.logLevel = lvl
}

//
// Tools - The tools in this section should only be used in an injector or for
// testing purposes.
//

// ConfigFile load di from configuration file
func (s *DI) ConfigFile() *config.File {
	if s.cacheConfigFile != nil {
		return s.cacheConfigFile
	}
	conf := &config.File{
		PathUserConf:     s.PathUserConf,
		PathTemplateConf: s.PathTemplateConf,
		Stderr:           s.Stderr,
	}
	s.validate(conf)
	err := conf.Load()
	errs.Panic(err)
	s.cacheConfigFile = conf
	return conf
}

// validate validate objects. Panics on failure
func (s *DI) validate(item interface{}) {
	validate := s.MakeValidate()
	err := validate.Struct(item)
	errs.PanicF("Validation Error: %v", err)
}

//
// Dependency Injectors
//

// MakeORM return ORM object.
func (s *DI) MakeORM() *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)
	if s.cacheORM != nil {
		return s.cacheORM
	}
	if _, err = os.Stat(s.SqliteDB); err == nil {
		db, err = gorm.Open(sqlite.Open(s.SqliteDB), &gorm.Config{
			NamingStrategy: &schema.NamingStrategy{
				SingularTable: true,
			},
		})
		errs.FatalF("Can not open ORM database %s: %v", s.SqliteDB, err)

	}
	s.cacheORM = db
	return db
}

// MakeLoggerOutput inject logger.OutputIface.
func (s *DI) MakeLoggerOutput(prefix string) *logger.Router {
	toFile := logger.New(s.MakeLogFile(), prefix, log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix, s.logLevel, s.MakeExitHandler())
	toStderr := logger.New(s.Stderr, prefix, log.Lmsgprefix, s.logLevel, s.MakeExitHandler())
	return logger.NewRouter(toFile, toStderr)
}

// MakeErrorHandler dependency injector
func (s *DI) MakeErrorHandler() *errs.Handler {
	if s.cacheErrHandler != nil {
		return s.cacheErrHandler
	}
	handler := errs.New(s.MakeExitHandler(), s.MakeLoggerOutput(""))
	s.cacheErrHandler = handler
	return handler
}

// MakeExitHandler dependency injector
func (s *DI) MakeExitHandler() *exit.Handler {
	if s.cacheExitHandler != nil {
		return s.cacheExitHandler
	}
	handler := &exit.Handler{
		Mode: s.ExitMode,
	}
	s.validate(handler)
	s.cacheExitHandler = handler
	return handler
}

// MakeCheck dependency injector
func (s *DI) MakeCheck() *check.Check {
	if s.cacheCheck != nil {
		return s.cacheCheck
	}
	chk := &check.Check{
		ConfigPath:         s.PathUserConf,
		ConfigTemplatePath: s.PathTemplateConf,
		HomeDir:            s.Home,
		MetadataDir:        s.PathMetadataDir,
		SQLiteFile:         s.SqliteDB,
		Stderr:             s.Stderr,
		Stdout:             s.Stdout,
		TemplateDir:        s.PathTemplateDir,
	}
	s.validate(chk)
	s.cacheCheck = chk
	return chk
}

// MakeSetup dependency injector
func (s *DI) MakeSetup() *setup.Setup {
	if s.cacheSetup != nil {
		return s.cacheSetup
	}
	i := &setup.Setup{
		ConfigPath:         s.PathUserConf,
		ConfigTemplatePath: s.PathTemplateConf,
		HomeDir:            s.Home,
		MetadataDir:        s.PathMetadataDir,
		SQLiteFile:         s.SqliteDB,
		TemplateDir:        s.PathTemplateDir,
	}
	s.validate(i)
	s.cacheSetup = i
	return i
}

// MakeInit dependency injector
func (s *DI) MakeInit() *initialize.Init {
	if s.cacheInit != nil {
		return s.cacheInit
	}
	i := &initialize.Init{
		ErrHandler: s.MakeErrorHandler(),
		Log:        s.MakeLoggerOutput(""),
	}
	s.validate(i)
	s.cacheInit = i
	return i
}

// MakeInstall dependency injector
func (s *DI) MakeInstall() *install.Install {
	if s.cacheInstall != nil {
		return s.cacheInstall
	}
	i := &install.Install{
		ConfigFile: s.ConfigFile(),
		ORM:        s.MakeORM(),
		Err:        s.MakeErrorHandler(),
		Log:        s.MakeLoggerOutput(""),
		Plumb:      s.MakePlumbingTemplate(),
		Stderr:     s.Stderr,
		Stdin:      s.Stdin,
		Stdout:     s.Stdout,
		Sync:       s.MakeSync(),
	}
	s.cacheInstall = i
	return i
}

// MakeList dependency injector
func (s *DI) MakeList() *list.List {
	if s.cacheList != nil {
		return s.cacheList
	}
	l := &list.List{
		Stderr: s.Stderr,
		Stdout: s.Stdout,
		Conf:   s.ConfigFile(),
	}
	s.validate(l)
	s.cacheList = l
	return l
}

// MakeLogFile create a logfile and return the interface
func (s *DI) MakeLogFile() *os.File {
	if s.cacheLogFile != nil {
		return s.cacheLogFile
	}
	var (
		tmpDir string
		f      *os.File
		err    error
	)
	if runtime.GOOS == "darwin" {
		tmpDir = "/tmp"
	} else {
		tmpDir = os.TempDir()
	}
	logPath := filepath.Join(tmpDir, "kick.log")

	fInfo, err := os.Stat(logPath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	} else if err == nil && fInfo.Size() > 1024*1024*2 {
		// Remove files greater than 2M
		os.Remove(logPath)
	}
	f, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	s.cacheLogFile = f
	return f
}

// MakePlumbingRepo injects di for plumbing.Plumb
func (s *DI) MakePlumbingRepo() *plumbing.Plumbing {
	if s.cachePlumbingRepo != nil {
		return s.cachePlumbingRepo
	}
	p := &plumbing.Plumbing{
		Base: s.PathRepoDir,
	}
	s.cachePlumbingRepo = p
	return p
}

// MakePlumbingTemplate injects di for plumbing.Plumb
func (s *DI) MakePlumbingTemplate() *plumbing.Plumbing {
	if s.cachePlumbingTemplate != nil {
		return s.cachePlumbingTemplate
	}
	p := &plumbing.Plumbing{
		Base: s.PathTemplateDir,
	}
	s.cachePlumbingTemplate = p
	return p
}

// MakeRemove dependency injector
func (s *DI) MakeRemove() *remove.Remove {
	if s.cacheRemove != nil {
		return s.cacheRemove
	}
	r := &remove.Remove{
		Conf:             s.ConfigFile(),
		Err:              s.MakeErrorHandler(),
		Log:              s.MakeLoggerOutput(""),
		PathTemplateConf: s.PathTemplateConf,
		PathUserConf:     s.PathUserConf,
		Stderr:           s.Stderr,
		Stdout:           s.Stdout,
	}
	s.cacheRemove = r
	return r
}

// MakeRepo dependency injector
func (s *DI) MakeRepo() *repo.Repo {
	wd, err := os.Getwd()
	errs.Panic(err)
	r := &repo.Repo{
		WD:         wd,
		Plumb:      s.MakePlumbingTemplate(),
		Validate:   s.MakeValidate(),
		ErrHandler: s.MakeErrorHandler(),
		Log:        s.MakeLoggerOutput(""),
	}
	return r
}

// MakeSearch dependency injector
func (s *DI) MakeSearch() *search.Search {
	if s.cacheSearch != nil {
		return s.cacheSearch
	}
	srch := &search.Search{
		ORM:    s.MakeORM(),
		Writer: os.Stdout,
	}
	s.cacheSearch = srch
	return srch
}

// MakeSync dependency injector
func (s *DI) MakeSync() *sync.Sync {
	if s.cacheSync != nil {
		return s.cacheSync
	}
	syn := &sync.Sync{
		ORM:                s.MakeORM(),
		Config:             s.ConfigFile(),
		ConfigTemplatePath: s.PathTemplateConf,
		Log:                s.MakeLoggerOutput(""),
		PlumbRepo:          s.MakePlumbingRepo(),
		PlumbTemplates:     s.MakePlumbingTemplate(),
		Stderr:             s.Stderr,
		Stdout:             s.Stdout,
	}
	s.cacheSync = syn
	return syn
}

// MakeTemplate dependency injector
func (s *DI) MakeTemplate() *template.Template {
	if s.cacheTemplate != nil {
		return s.cacheTemplate
	}
	vars := variables.New()
	vars.ProjectVariable("NAME", s.ProjectName)
	t := &template.Template{
		Config:        s.ConfigFile(),
		Log:           s.MakeLoggerOutput(""),
		Errs:          s.MakeErrorHandler(),
		Exit:          s.MakeExitHandler(),
		Stderr:        s.Stderr,
		Stdout:        s.Stdout,
		TemplateDir:   s.PathTemplateDir,
		Variables:     vars,
		RenderCurrent: "envsubst",
		RenderersAvail: map[string]renderer.Renderer{
			"texttemplate": &renderer.RenderText{},
			"envsubst":     &renderer.RenderEnv{},
		},
	}
	s.cacheTemplate = t
	return t
}

// MakeUpdate dependency injector
func (s *DI) MakeUpdate() *update.Update {
	if s.cacheUpdate != nil {
		return s.cacheUpdate
	}
	u := &update.Update{
		ConfigFile:  s.ConfigFile(),
		ORM:         s.MakeORM(),
		Log:         s.MakeLoggerOutput(""),
		MetadataDir: s.PathMetadataDir,
	}
	s.cacheUpdate = u
	return u
}

// MakeValidate dependency injector
func (s *DI) MakeValidate() *validator.Validate {
	v := validator.New()
	return v
}
