package iinitialize

import (
	"github.com/kick-project/kick/internal/di"
)

// Inject creates di for initialize.New
func Inject(s *di.DI) (opts struct {
	ConfigPath         string // Path to configuration file
	ConfigTemplatePath string // Path to template configuration file
	HomeDir            string // Path to home directory
	MetadataDir        string // Path to metadata directory
	SQLiteFile         string // Path to DB file
	TemplateDir        string // Path to template directory
}) {
	opts.ConfigPath = s.PathUserConf
	opts.ConfigTemplatePath = s.PathTemplateConf
	opts.HomeDir = s.Home
	opts.MetadataDir = s.PathMetadataDir
	opts.SQLiteFile = s.SqliteDB
	opts.TemplateDir = s.PathTemplateDir
	return opts
}
