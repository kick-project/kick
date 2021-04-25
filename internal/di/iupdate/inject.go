package iupdate

import (
	"github.com/apex/log"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/config"
	"gorm.io/gorm"
)

// Inject creates di for metadata.New
func Inject(s *di.DI) (opts struct {
	ConfigFile  *config.File
	ORM         *gorm.DB
	Log         *log.Logger
	MetadataDir string
}) {
	opts.ConfigFile = s.ConfigFile()
	opts.ORM = s.GetORM()
	opts.Log = s.GetLogger()
	opts.MetadataDir = s.PathMetadataDir
	return opts
}
