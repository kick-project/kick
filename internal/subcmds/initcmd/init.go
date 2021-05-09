package initcmd

import (
	"log"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `initialize configuration

Usage:
    kick init

Options:
    -h --help     print help
`

// OptInit initialize configuration file
type OptInit struct {
	Init bool `docopt:"init"`
}

// InitCmd initialize configuration
func InitCmd(args []string, inject *di.DI) int {
	opts := &OptInit{}
	options.Bind(UsageDoc, args, opts)
	if !opts.Init {
		log.Println("error can not initialize")
		return 256
	}

	i := inject.MakeInitialize()
	i.Init()

	return 0
}
