package itablesync

import (
	"database/sql"

	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/settings"
)

// Inject creates settings for tablesync.New
func Inject(s *settings.Settings) (opts struct {
	DB         *sql.DB
	Config     *config.File
	ConfigPath string
}) {
	opts.DB = s.GetDB()
	opts.Config = s.ConfigFile()
	opts.ConfigPath = s.PathUserConf

	return opts
}
