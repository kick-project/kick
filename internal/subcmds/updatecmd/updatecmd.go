package updatecmd

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/services/update"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/icheck"
	"github.com/kick-project/kick/internal/settings/iupdate"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `update repository data

Usage:
    kick update

Options:
    -h --help     print help
`

// OptUpdate bindings for docopts
type OptUpdate struct {
	Search bool `docopt:"update"`
}

// Update for templates
func Update(args []string, s *settings.Settings) int {
	opts := &OptUpdate{}
	options.Bind(usageDoc, args, opts)

	chk := &check.Check{}
	err := copier.Copy(chk, icheck.Inject(s))
	errutils.Epanic(err)

	if err = chk.Init(); err != nil {
		fmt.Fprintf(s.Stderr, "%s\n", err.Error())
		utils.Exit(255)
	}

	u := &update.Update{}
	err = copier.Copy(u, iupdate.Inject(s))
	errutils.Epanic(err)
	err = u.Build()
	errutils.Epanic(err)

	return 0
}
