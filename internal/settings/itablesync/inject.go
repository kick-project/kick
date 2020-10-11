package itablesync

import (
	"database/sql"

	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/settings"
)

// Inject creates settings for tablesync.New
func Inject(s *settings.Settings) (opts struct {
	DB                 *sql.DB
	Config             *config.File
	ConfigTemplatePath string
}) {
	opts.DB = s.GetDB()
	opts.Config = s.ConfigFile()
	opts.ConfigTemplatePath = s.PathTemplateConf

	return opts
}
