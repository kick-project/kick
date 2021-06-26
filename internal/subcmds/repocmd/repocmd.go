package repocmd

import (
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `Buid/list/inform on repositories WIP

Usage:
    kick repo build
    kick repo list
    kick repo info <repo>

Options:
    -h --help    print help
    repo         repo subcommand
    build        build repo by downloading the URLS defined in repo.yml and creating the files templates/*.yml
    list         list repositories
    info         repository and/or template information
    <repo>       name of repository
`

// OptRepo initialize configuration file
type OptRepo struct {
	Repo     bool   `docopt:"repo"`
	Build    bool   `docopt:"build"`
	List     bool   `docopt:"list"`
	Info     bool   `docopt:"info"`
	RepoName string `docopt:"<repo>"`
}

// Repo install a template
func Repo(args []string, inject *di.DI) int {
	opts := &OptRepo{}
	options.Bind(UsageDoc, args, opts)
	r := inject.MakeRepo()
	switch {
	case opts.Build:
		r.Build()
	case opts.List:
		r.List()
	case opts.Info:
		r.Info(opts.RepoName)
	}
	return 0
}
