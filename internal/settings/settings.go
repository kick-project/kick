package settings

import (
	"database/sql"
	"fmt"
	"os"
	fp "path/filepath"

	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/services/template"
	"github.com/crosseyed/prjstart/internal/utils/dfaults"
	_ "github.com/mattn/go-sqlite3" // Driver for database/sql
)

//
// Settings
//

// Settings provides settings for resources & services.
type Settings struct {
	home         string
	File         *config.File
	DBConfig     *DBConfig
	TemplateOpts *template.Options
}

// GetSettings get settings
func GetSettings(home string) *Settings {
	s := &Settings{
		home: dfaults.String(homePath(), home),
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
		Path: fp.Clean(fmt.Sprintf("%s/.prjstart.yml", s.home)),
	})
	conf.Load()
	s.File = conf
	return s.File
}

// DBConf configure database settings
func (s *Settings) DBConf() *DBConfig {
	if s.DBConfig != nil {
		return s.DBConfig
	}
	path := fp.Clean(fmt.Sprintf("%s/.prjstart/metadata/metadata.db", s.home))
	db := &DBConfig{
		SQLitePath: path,
		Driver:     "sqlite3",
		DSN:        fmt.Sprintf("file:%s?_foreign_key=on", path),
	}
	s.DBConfig = db
	return s.DBConfig
}

// Template injects template settings
func (s *Settings) Template() *template.Options {
	if s.TemplateOpts != nil {
		return s.TemplateOpts
	}
	configFile := s.ConfigFile()
	templateDir := fp.Clean(fmt.Sprintf("%s/.prjstart/templates", s.home))
	vars := template.NewTmplVars()
	o := &template.Options{
		Config:      configFile,
		TemplateDir: templateDir,
		Variables:   vars,
	}
	s.TemplateOpts = o

	return s.TemplateOpts
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

//
// Misc
//

func homePath() string {
	return os.Getenv("HOME")
}
