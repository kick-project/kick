package internal

import (
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/docopt/docopt-go"
)

//
// Document strings
//

var UsageMain = `Generate project scaffolding from a predefined set of templates

Usage:
    prjstart start
    prjstart list

Options:
    -h --help     Print help.
    -v --version  Print version.
    start         Start a project.
    list          List projects/variables.
`
var UsageMainRemoteFeature = `Generate project scaffolding from a predefined set of templates

Usage:
    prjstart start
    prjstart list
    prjstart install

Options:
    -h --help     Print help.
    -v --version  Print version.
    start         Start a project.
    list          List available project options.
    install       Install a template.
`

var UsageStart = `Generate project scaffolding

Usage:
    prjstart start <template> <project>

Options:
    -h --help     Print help.
    <template>    Template name.
    <project>     Project name.
`

var UsageList = `List templates

Usage:
    prjstart list [--url]
    prjstart list [--vars]

Options:
    -h --help     Print help.
    -u --url      Print URL.
    -v --vars     Show Variables.
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
	Install bool `docopt:"install"`
}

func GetOptMain(args []string) *OptMain {
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
	opts, err := docopt.ParseArgs(UsageMain, filterArgs, Version)
	utils.ChkErr(err, utils.Epanicf, "Can not parse usage doc: %s", err) // nolint
	o := new(OptMain)
	err = opts.Bind(o)
	utils.ChkErr(err, utils.Epanicf, "Can not bind to structure: %s", err) // nolint
	return o
}

type OptStart struct {
	Start   bool   `docopt:"start"`
	Tmpl    string `docopt:"<template>"`
	Project string `docopt:"<project>"`
}

func GetOptStart(args []string) *OptStart {
	opts, err := docopt.ParseArgs(UsageStart, args, "")
	utils.ChkErr(err, utils.Epanicf, "Can not parse usage doc: %s", err) // nolint
	o := new(OptStart)
	err = opts.Bind(o)
	utils.ChkErr(err, utils.Epanicf, "Can not bind to structure: %s", err) // nolint
	return o
}

type OptList struct {
	List   bool `docopt:"list"`
	Local  bool `docopt:"--local"`
	Remote bool `docopt:"--remote"`
	All    bool `docopt:"--all"`
	Url    bool `docopt:"--url"`
	Vars   bool `docopt:"--vars"`
}

func GetOptList(args []string) *OptList {
	opts, err := docopt.ParseArgs(UsageList, args, "")
	utils.ChkErr(err, utils.Epanicf, "Can not parse usage doc: %s", err) // nolint
	o := new(OptList)
	err = opts.Bind(o)
	utils.ChkErr(err, utils.Epanicf, "Can not bind to structure: %s", err) // nolint
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
	utils.ChkErr(err, utils.Epanicf, "Can not parse usage doc: %s", err) // nolint
	o := new(OptInstall)
	err = opts.Bind(o)
	utils.ChkErr(err, utils.Epanicf, "Can not bind to structure: %s", err) // nolint
	return o
}
