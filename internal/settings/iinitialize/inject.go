package iinitialize

import (
	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/settings"
)

// Inject creates settings for initialize.New
func Inject(s *settings.Settings) (opts struct {
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
	opts.DBDriver = s.DBDriver
	opts.DSN = s.DBDsn
	opts.HomeDir = s.Home
	opts.MetadataDir = s.PathMetadataDir
	opts.SQLiteFile = s.SqliteDB
	opts.TemplateDir = s.PathTemplateDir
	return opts
}
