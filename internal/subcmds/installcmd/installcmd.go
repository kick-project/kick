package installcmd

import (
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/icheck"
	"github.com/kick-project/kick/internal/di/iinstall"
	"github.com/kick-project/kick/internal/di/iupdate"
	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/services/install"
	"github.com/kick-project/kick/internal/services/update"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `Install template

Usage:
    kick install <handle> <location>

Options:
    -h --help        print help
    <handle>         name to use when creating new projects
    <location>       template name, URL or location of template
`

// OptInstall initialize configuration file
type OptInstall struct {
	Install  bool   `docopt:"install"`
	Template string `docopt:"<location>"`
	Handle   string `docopt:"<handle>"`
}

// Install install a template
func Install(args []string, inject *di.DI) int {
	opts := &OptInstall{}
	options.Bind(UsageDoc, args, opts)
	if !opts.Install {
		errutils.Epanic(errors.New("Install set to false"))
		return 256
	}

	chk := &check.Check{}
	err := copier.Copy(chk, icheck.Inject(inject))
	errutils.Epanic(err)

	if err = chk.Init(); err != nil {
		fmt.Fprintf(inject.Stderr, "%s\n", err.Error())
		exit.Exit(255)
	}

	m := &update.Update{}
	err = copier.Copy(m, iupdate.Inject(inject))
	errutils.Epanic(err)

	err = m.Build()
	errutils.Epanic(err)

	inst := &install.Install{}
	err = copier.Copy(inst, iinstall.Inject(inject))
	errutils.Epanic(err)
	return inst.Install(opts.Handle, opts.Template)
}
