package iupdate

import (
	"database/sql"

	"github.com/apex/log"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/settings"
)

// Inject creates settings for metadata.New
func Inject(s *settings.Settings) (opts struct {
	ConfigFile  *config.File
	DB          *sql.DB
	Log         *log.Logger
	MetadataDir string
}) {
	db := s.GetDB()
	opts.ConfigFile = s.ConfigFile()
	opts.DB = db
	opts.Log = s.GetLogger()
	opts.MetadataDir = s.PathMetadataDir
	return opts
}
