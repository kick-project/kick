package removecmd

import (
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/services/remove"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/icheck"
	"github.com/kick-project/kick/internal/settings/iremove"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `Remove an installed template

Usage:
    kick remove <handle>

Options:
    -h --help        print help
    <handle>         handle to remove
`

// OptRemove remove an installed template
type OptRemove struct {
	Remove bool   `docopt:"remove"`
	Handle string `docopt:"<handle>"`
}

// Remove remove a template
func Remove(args []string, s *settings.Settings) int {
	opts := &OptRemove{}
	options.Bind(usageDoc, args, opts)
	if !opts.Remove {
		errutils.Epanic(errors.New("Remove set to false"))
		return 256
	}

	chk := &check.Check{}
	err := copier.Copy(chk, icheck.Inject(s))
	errutils.Epanic(err)
	if err = chk.Init(); err != nil {
		fmt.Fprintf(s.Stderr, "%s\n", err.Error())
		utils.Exit(255)
	}

	rm := &remove.Remove{}
	err = copier.Copy(rm, iremove.Inject(s))
	errutils.Epanic(err)
	return rm.Remove(opts.Handle)
}
