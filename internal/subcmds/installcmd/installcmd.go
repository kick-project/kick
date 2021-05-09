package installcmd

import (
	"errors"
	"fmt"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/options"
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
		errs.Panic(errors.New("Install set to false"))
		return 256
	}

	chk := inject.MakeCheck()

	if err := chk.Init(); err != nil {
		fmt.Fprintf(inject.Stderr, "%s\n", err.Error())
		exit.Exit(255)
	}

	m := inject.MakeUpdate()
	err := m.Build()
	errs.Panic(err)

	inst := inject.MakeInstall()
	return inst.Install(opts.Handle, opts.Template)
}
