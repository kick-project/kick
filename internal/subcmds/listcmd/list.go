package listcmd

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/services/list"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/icheck"
	"github.com/kick-project/kick/internal/settings/ilist"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `List handles/templates

Usage:
    kick list [-l]

Options:
	-h --help     Print help.
    -l            Print Long output.
`

type OptList struct {
	List bool `docopt:"list"`
	Long bool `docopt:"-l"`
}

// List starts the list sub command
func List(args []string, s *settings.Settings) int {
	opts := &OptList{}
	options.Bind(usageDoc, args, opts)

	chk := &check.Check{}
	err := copier.Copy(chk, icheck.Inject(s))
	errutils.Epanic(err)
	if err = chk.Init(); err != nil {
		fmt.Fprintf(s.Stderr, "%s\n", err.Error())
		utils.Exit(255)
	}

	l := &list.List{}
	err = copier.Copy(l, ilist.Inject(s))
	errutils.Epanic(err)

	return l.List(opts.Long)
}
