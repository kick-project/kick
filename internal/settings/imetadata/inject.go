package imetadata

import (
	"database/sql"

	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/settings"
)

// Inject creates settings for metadata.New
func Inject(s *settings.Settings) (opts struct {
	ConfigFile  *config.File
	MetadataDir string
	DB          *sql.DB
}) {
	db := s.GetDB()
	opts.ConfigFile = s.ConfigFile()
	opts.MetadataDir = s.PathMetadataDir
	opts.DB = db
	return opts
}
