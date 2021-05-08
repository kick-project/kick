package listcmd

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/icheck"
	"github.com/kick-project/kick/internal/di/ilist"
	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/services/list"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `list handles/templates

Usage:
    kick list [-l]

Options:
    -h --help     print help
    -l            print long output
`

// OptList docopts options to list installed templates
type OptList struct {
	List bool `docopt:"list"`
	Long bool `docopt:"-l"`
}

// List starts the list sub command
func List(args []string, inject *di.DI) int {
	opts := &OptList{}
	options.Bind(usageDoc, args, opts)

	chk := &check.Check{}
	err := copier.Copy(chk, icheck.Inject(inject))
	errutils.Epanic(err)
	if err = chk.Init(); err != nil {
		fmt.Fprintf(inject.Stderr, "%s\n", err.Error())
		exit.Exit(255)
	}

	l := &list.List{}
	err = copier.Copy(l, ilist.Inject(inject))
	errutils.Epanic(err)

	return l.List(opts.Long)
}
