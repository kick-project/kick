package installcmd

import (
	"errors"

	"github.com/crosseyed/prjstart/internal/services/install"
	"github.com/crosseyed/prjstart/internal/services/metadata"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/settings/iinstall"
	"github.com/crosseyed/prjstart/internal/settings/imetadata"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/crosseyed/prjstart/internal/utils/options"
	"github.com/jinzhu/copier"
)

var usageDoc = `Install template

Usage:
    prjstart install <handle> <location>

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
	if opts.Install == false {
		errutils.Epanic(errors.New("Install set to false"))
		return 256
	}

	m := &metadata.Metadata{}
	copier.Copy(m, imetadata.Inject(s))
	m.Build()

	inst := &install.Install{}
	copier.Copy(inst, iinstall.Inject(s))
	return inst.Install(opts.Handle, opts.Template)
}
