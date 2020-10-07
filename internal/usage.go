package internal

import (
	"github.com/crosseyed/prjstart/internal/fflags"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/docopt/docopt-go"
)

//
// Document strings
//

var UsageMain = `Generate project scaffolding from a predefined set of templates

Usage:
    prjstart start
    prjstart list
    prjstart init

Options:
    -h --help     Print help.
    -v --version  Print version.
    start         Start a project.
    list          List projects/variables.
    init          Initialize configuration.
`
var UsageMainFFREMOTE = `Generate project scaffolding from a predefined set of templates

Usage:
    prjstart start
    prjstart list
    prjstart init
    prjstart search
    prjstart install

Options:
    -h --help     Print help.
    -v --version  Print version.
    start         Start a project.
    list          List available project options.
    init          Initialize configuration.
    search        Search for available templates.
    install       Install a template.
`

var UsageSearch = `Search for templates

Usage:
    prjstart search [--long] <template>

Options:
    -h --help     Print help.
    -l --long     Long output.
`

var UsageListRemoteFeature = `List templates

Usage:
    prjstart list [--local|--remote|--all] [--url]
    prjstart list [--vars]

Options:
    -h --help     Print help.
    -u --url      Print URL.
    -v --vars     Show Variables.
    -l --local    List local templates. [default]
    -r --remote   List remote templates listed in sets.
    -a --all      List local and remote set templates
`

var UsageInstall = `Install a remote set

Usage:
    prjstart install <set> <template> [<name>]

Options:
    -h --help     Print help.
    <set>         Install from set
    <template>    Install template
    <name>        Optional name
`

//
// Options
//

// OptMain loads the data parsed from the root "prjstart" command
type OptMain struct {
	Start   bool `docopt:"start"`
	List    bool `docopt:"list"`
	Init    bool `docopt:"init"`
	Search  bool `docopt:"search"`
	Install bool `docopt:"install"`
}

func GetOptMain(args []string) *OptMain {
	var (
		opts docopt.Opts
		err  error
	)
	filterArgs := []string{}
	i := 0
	for _, arg := range args {
		i++
		if i == 1 {
			continue
		}
		filterArgs = append(filterArgs, arg)
		break
	}
	if fflags.Remote() {
		opts, err = docopt.ParseArgs(UsageMainFFREMOTE, filterArgs, Version)
	} else {
		opts, err = docopt.ParseArgs(UsageMain, filterArgs, Version)
	}
	errutils.Epanicf("Can not parse usage doc: %s", err) // nolint
	o := new(OptMain)
	err = opts.Bind(o)
	errutils.Epanicf("Can not bind to structure: %s", err) // nolint
	return o
}

type OptSearch struct {
	Search   bool   `docopt:"search"`
	Long     bool   `docopt:"--long"`
	Template string `docopt:"<template>"`
}

func GetOptSearch(args []string) *OptSearch {
	opts, err := docopt.ParseArgs(UsageSearch, args, "")
	errutils.Epanicf("Can not parse usage doc: %s", err) // nolint
	o := new(OptSearch)
	err = opts.Bind(o)
	errutils.Epanicf("Can not bind to structure: %s", err) // nolint
	return o
}

type OptInstall struct {
	Install  bool   `docopt:"install"`
	Set      string `docopt:"<set>"`
	Template string `docopt:"<template>"`
	Name     string `docopt:"<name>"`
}

func GetOptInstall(args []string) *OptInstall {
	opts, err := docopt.ParseArgs(UsageInstall, args, "")
	errutils.Epanicf("Can not parse usage doc: %s", err) // nolint
	o := new(OptInstall)
	err = opts.Bind(o)
	errutils.Epanicf("Can not bind to structure: %s", err) // nolint
	return o
}
