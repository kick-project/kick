package setupcmd

import (
	"log"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `initialize configuration

Usage:
    kick setup

Options:
    -h --help     print help
`

// OptSetup initialize configuration file
type OptSetup struct {
	Setup bool `docopt:"setup"`
}

// SetupCmd initialize configuration
func SetupCmd(args []string, inject *di.DI) int {
	opts := &OptSetup{}
	options.Bind(UsageDoc, args, opts)
	if !opts.Setup {
		log.Println("error can not initialize")
		return 256
	}

	i := inject.MakeSetup()
	i.Init()

	return 0
}
