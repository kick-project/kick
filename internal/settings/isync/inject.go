package isync

import (
	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/iplumbing"
	"github.com/kick-project/kick/internal/utils/errutils"
	"gorm.io/gorm"
)

// Inject creates settings for tablesync.New
func Inject(s *settings.Settings) (opts struct {
	BasePath           string
	ORM                *gorm.DB
	Config             *config.File
	ConfigTemplatePath string
	Plumb              *plumbing.Plumbing
}) {
	plumb := &plumbing.Plumbing{}
	err := copier.Copy(plumb, iplumbing.Inject(s))
	errutils.Epanic(err)
	opts.BasePath = s.PathDownloadDir
	opts.ORM = s.GetORM()
	opts.Config = s.ConfigFile()
	opts.ConfigTemplatePath = s.PathTemplateConf
	opts.Plumb = plumb
	return opts
}
