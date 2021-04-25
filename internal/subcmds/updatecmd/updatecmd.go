package updatecmd

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/icheck"
	"github.com/kick-project/kick/internal/di/iupdate"
	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/services/update"
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
func Update(args []string, inject *di.DI) int {
	opts := &OptUpdate{}
	options.Bind(usageDoc, args, opts)

	chk := &check.Check{}
	err := copier.Copy(chk, icheck.Inject(inject))
	errutils.Epanic(err)

	if err = chk.Init(); err != nil {
		fmt.Fprintf(inject.Stderr, "%s\n", err.Error())
		utils.Exit(255)
	}

	u := &update.Update{}
	err = copier.Copy(u, iupdate.Inject(inject))
	errutils.Epanic(err)
	err = u.Build()
	errutils.Epanic(err)

	return 0
}
