package iinstall

import (
	"database/sql"
	"io"

	"github.com/apex/log"
	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/resources/gitclient/plumbing"
	"github.com/crosseyed/prjstart/internal/resources/sync"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/settings/iplumbing"
	"github.com/crosseyed/prjstart/internal/settings/isync"
	"github.com/jinzhu/copier"
)

// Inject inject options for install.Install
func Inject(s *settings.Settings) (opts struct {
	ConfigFile *config.File
	DB         *sql.DB
	Log        *log.Logger
	Plumb      *plumbing.Plumbing
	Stderr     io.Writer
	Stdin      io.Reader
	Stdout     io.Writer
	Sync       *sync.Sync
}) {
	synchro := &sync.Sync{}
	copier.Copy(synchro, isync.Inject(s))
	plumb := &plumbing.Plumbing{}
	copier.Copy(plumb, iplumbing.Inject(s))
	opts.ConfigFile = s.ConfigFile()
	opts.DB = s.GetDB()
	opts.Log = s.GetLogger()
	opts.Plumb = plumb
	opts.Stderr = s.Stderr
	opts.Stdin = s.Stdin
	opts.Stdout = s.Stdout
	opts.Sync = synchro
	return
}
