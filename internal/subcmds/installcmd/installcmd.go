package installcmd

import (
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/services/install"
	"github.com/kick-project/kick/internal/services/update"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/icheck"
	"github.com/kick-project/kick/internal/settings/iinstall"
	"github.com/kick-project/kick/internal/settings/iupdate"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `Install template

Usage:
    kick install <handle> <location>

Options:
    -h --help        Print help.
    <handle>         Name to use when creating new projects
    <location>       Template name, URL or location of template.
`

// OptInstall initialize configuration file
type OptInstall struct {
	Install  bool   `docopt:"install"`
	Template string `docopt:"<location>"`
	Handle   string `docopt:"<handle>"`
}

// Install install a template
func Install(args []string, s *settings.Settings) int {
	opts := &OptInstall{}
	options.Bind(usageDoc, args, opts)
	if !opts.Install {
		errutils.Epanic(errors.New("Install set to false"))
		return 256
	}

	chk := &check.Check{}
	err := copier.Copy(chk, icheck.Inject(s))
	errutils.Epanic(err)

	if err = chk.Init(); err != nil {
		fmt.Fprintf(s.Stderr, "%s\n", err.Error())
		utils.Exit(255)
	}

	m := &update.Update{}
	err = copier.Copy(m, iupdate.Inject(s))
	errutils.Epanic(err)

	err = m.Build()
	errutils.Epanic(err)

	inst := &install.Install{}
	err = copier.Copy(inst, iinstall.Inject(s))
	errutils.Epanic(err)
	return inst.Install(opts.Handle, opts.Template)
}
