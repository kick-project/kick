package icheck

import (
	"io"

	"github.com/kick-project/kick/internal/settings"
)

// Inject creates settings for check.New
func Inject(s *settings.Settings) (opts struct {
	ConfigPath         string
	ConfigTemplatePath string
	HomeDir            string
	MetadataDir        string
	SQLiteFile         string
	Stderr             io.Writer
	Stdout             io.Writer
	TemplateDir        string
}) {
	opts.ConfigPath = s.PathUserConf
	opts.ConfigTemplatePath = s.PathTemplateConf
	opts.HomeDir = s.Home
	opts.MetadataDir = s.PathMetadataDir
	opts.SQLiteFile = s.SqliteDB
	opts.Stderr = s.Stderr
	opts.Stdout = s.Stdout
	opts.TemplateDir = s.PathTemplateDir
	return opts
}
