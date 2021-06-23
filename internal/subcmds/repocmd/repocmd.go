package repocmd

import (
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/fflags"
	"github.com/kick-project/kick/internal/resources/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `Build a repo from repo.yml

Usage:
    kick repo build

Options:
    -h --help    print help
    repo         repo subcommand
    build        build repo by downloading the URLS defined in repo.yml and creating the files templates/*.yml
`

// UsageDocExtended help document passed to docopts
var UsageDocExtended = `Build a repo from repo.yml

Usage:
    kick repo build
	kick repo list
	kick repo info [--tver] <repo> [<template>]

Options:
    -h --help    print help
	--tver       template versions  
    repo         repo subcommand
    build        build repo by downloading the URLS defined in repo.yml and creating the files templates/*.yml
	list         list repositories
	info         repository and/or template information
`

// OptRepo initialize configuration file
type OptRepo struct {
	Repo         bool   `docopt:"repo"`
	Build        bool   `docopt:"build"`
	List         bool   `docopt:"list"`
	Info         bool   `docopt:"info"`
	TVer         bool   `docopt:"--tver"`
	RepoName     string `docopt:"<repo>"`
	TemplateName string `docopt:"<template>"`
}

// Repo install a template
func Repo(args []string, inject *di.DI) int {
	if fflags.RepoExtension() {
		return repoExtension(args, inject)
	}

	return repoOrig(args, inject)
}

func repoOrig(args []string, inject *di.DI) int {
	opts := &OptRepo{}
	options.Bind(UsageDoc, args, opts)
	if opts.Build {
		r := inject.MakeRepo()
		r.Build()
	}
	return 255
}

func repoExtension(args []string, inject *di.DI) int {
	opts := &OptRepo{}
	options.Bind(UsageDoc, args, opts)
	r := inject.MakeRepo()
	switch {
	case opts.Build:
		r.Build()
	case opts.List:
	case opts.Info:
	}
	return 255
}
