package iremove

import (
	"database/sql"
	"io"

	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/settings"
)

// Inject inject settings for remove.Remove
func Inject(s *settings.Settings) (opts struct {
	Conf             *config.File
	DB               *sql.DB
	PathTemplateConf string
	PathUserConf     string
	Stderr           io.Writer
	Stdout           io.Writer
}) {
	opts.Conf = s.ConfigFile()
	opts.DB = s.GetDB()
	opts.PathTemplateConf = s.PathTemplateConf
	opts.PathUserConf = s.PathUserConf
	opts.Stderr = s.Stderr
	opts.Stdout = s.Stdout
	return
}
