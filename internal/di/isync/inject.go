package isync

import (
	"io"

	"github.com/apex/log"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/iplumbing"
	"github.com/kick-project/kick/internal/utils/errutils"
	"gorm.io/gorm"
)

// Inject creates di for tablesync.New
func Inject(s *di.DI) (opts struct {
	ORM                *gorm.DB
	Config             *config.File
	ConfigTemplatePath string
	Log                *log.Logger
	PlumbTemplates     *plumbing.Plumbing
	PlumbGlobal        *plumbing.Plumbing
	PlumbMaster        *plumbing.Plumbing
	Stderr             io.Writer
	Stdout             io.Writer
}) {
	plumbGlobal := &plumbing.Plumbing{}
	err := copier.Copy(plumbGlobal, iplumbing.InjectGlobal(s))
	errutils.Epanic(err)
	plumbMaster := &plumbing.Plumbing{}
	err = copier.Copy(plumbMaster, iplumbing.InjectMaster(s))
	errutils.Epanic(err)
	plumbTemplate := &plumbing.Plumbing{}
	err = copier.Copy(plumbTemplate, iplumbing.InjectTemplate(s))
	errutils.Epanic(err)

	opts.ORM = s.GetORM()
	opts.Config = s.ConfigFile()
	opts.ConfigTemplatePath = s.PathTemplateConf
	opts.Log = s.GetLogger()
	opts.PlumbGlobal = plumbGlobal
	opts.PlumbMaster = plumbMaster
	opts.PlumbTemplates = plumbTemplate
	opts.Stderr = s.Stderr
	opts.Stdout = s.Stdout
	return opts
}
