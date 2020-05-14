package subcmds

import (
	"github.com/crosseyed/prjstart/internal"
	"github.com/crosseyed/prjstart/internal/build"
	"github.com/crosseyed/prjstart/internal/globals"
)

func Start(args []string) int {
	opts := internal.GetOptStart(args)
	vars := internal.SetVars(opts)
	globals.Vars = vars.GetVars()

	g := build.Build{}
	g.SetSrc(opts.Tmpl)
	g.SetDest(opts.Project)
	ret := g.Run()
	return ret
}
