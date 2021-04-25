package initcmd

import (
	"log"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/iinitialize"
	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `initialize configuration

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
	options.Bind(usageDoc, args, opts)
	if !opts.Init {
		log.Println("error can not initialize")
		return 256
	}

	i := &initialize.Initialize{}
	err := copier.Copy(i, iinitialize.Inject(inject))
	errutils.Epanic(err)
	i.Init()

	return 0
}
