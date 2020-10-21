// kick:render
package options

import (
	"github.com/docopt/docopt-go"
)

func GetUsage(argv []string, version string) *Options {
	usage := `${PROJECT_NAME}

Usage:
  ${PROJECT_NAME}

Options:
  -h --help     Show this screen.
  --version     Show version.`

	opts, err := docopt.ParseArgs(usage, argv, version)
	if err != nil {
		panic(err)
	}
	config := &Options{}
	err = opts.Bind(config)
	if err != nil {
		panic(err)
	}
	return config
}

type Options struct {
}
