package iinstall

import (
	"io"

	"github.com/apex/log"
	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/iplumbing"
	"github.com/kick-project/kick/internal/di/isync"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/resources/sync"
	"github.com/kick-project/kick/internal/utils/errutils"
	"gorm.io/gorm"
)

// Inject inject options for install.Install
func Inject(s *di.DI) (opts struct {
	ConfigFile *config.File
	ORM        *gorm.DB
	Log        *log.Logger
	Plumb      *plumbing.Plumbing
	Stderr     io.Writer
	Stdin      io.Reader
	Stdout     io.Writer
	Sync       *sync.Sync
}) {
	synchro := &sync.Sync{}
	err := copier.Copy(synchro, isync.Inject(s))
	errutils.Epanic(err)
	plumb := &plumbing.Plumbing{}
	err = copier.Copy(plumb, iplumbing.InjectTemplate(s))
	errutils.Epanic(err)
	opts.ConfigFile = s.ConfigFile()
	opts.ORM = s.MakeORM()
	opts.Log = s.MakeLogger()
	opts.Plumb = plumb
	opts.Stderr = s.Stderr
	opts.Stdin = s.Stdin
	opts.Stdout = s.Stdout
	opts.Sync = synchro
	return
}
