package installcmd

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/services/install"
	"github.com/kick-project/kick/internal/services/update"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/iinstall"
	"github.com/kick-project/kick/internal/settings/iupdate"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `Install template

Usage:
    kick install <handle> <location>

Options:
    -h --help        Print help.
    <location>       Template name, URL or location of template
    <name>           Name to use when creating new projects
`

// OptInstall initialize configuration file
type OptInstall struct {
	Install  bool   `docopt:"install"`
	URL      string `docopt:"--url"`
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

	m := &update.Update{}
	err := copier.Copy(m, iupdate.Inject(s))
	errutils.Epanic(err)
	err = m.Build()
	errutils.Epanic(err)

	inst := &install.Install{}
	err = copier.Copy(inst, iinstall.Inject(s))
	errutils.Epanic(err)
	return inst.Install(opts.Handle, opts.Template)
}
