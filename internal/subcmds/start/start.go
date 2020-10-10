package start

import (
	"path/filepath"

	"github.com/crosseyed/prjstart/internal/resources/db/tablesync"
	"github.com/crosseyed/prjstart/internal/services/template"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/settings/itablesync"
	"github.com/crosseyed/prjstart/internal/settings/itemplate"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/docopt/docopt-go"
)

var usageDoc = `Generate project scaffolding

Usage:
    prjstart start <template> <project>

Options:
    -h --help     Print help.
    <template>    Template name.
    <project>     Project name.
`

type OptStart struct {
	Start   bool   `docopt:"start"`
	Tmpl    string `docopt:"<template>"`
	Project string `docopt:"<project>"`
}

// GetOptStart parse start options from document text
func GetOptStart(args []string) *OptStart {
	opts, err := docopt.ParseArgs(usageDoc, args, "")
	errutils.Epanicf("Can not parse usage doc: %s", err) // nolint
	o := new(OptStart)
	err = opts.Bind(o)
	errutils.Epanicf("Can not bind to structure: %s", err) // nolint
	return o
}

// Start start cli option
func Start(args []string, s *settings.Settings) int {
	opts := GetOptStart(args)

	// Sync DB table "installed" with configuration file
	sync := tablesync.New(itablesync.Inject(s))
	sync.SyncInstalled()

	templateOptions := itemplate.Inject(s)

	// Set project name
	name := filepath.Base(opts.Project)
	templateOptions.Variables.SetProjectVar("NAME", name)
	t := template.New(templateOptions)
	t.SetSrcDest(opts.Tmpl, opts.Project)
	ret := t.Run()
	return ret
}
