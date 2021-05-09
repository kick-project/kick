package repocmd

import (
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `Build a repo from repo.xml

Usage:
    kick repo build

Options:
    -h --help    print help
    repo         repo subcommand
    build        build repo by downloading the URLS defined in repo.xml and creating the files templates/*.yml
`

// OptRepo initialize configuration file
type OptRepo struct {
	Repo  bool `docopt:"repo"`
	Build bool `docopt:"build"`
}

// Repo install a template
func Repo(args []string, inject *di.DI) int {
	opts := &OptRepo{}
	options.Bind(UsageDoc, args, opts)
	if opts.Build {
		r := inject.MakeRepo()
		r.Build()
	}
	return 255
}
