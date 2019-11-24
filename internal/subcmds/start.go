package subcmds

import (
	"github.com/crosseyed/prjstart/internal"
)

func Start(args []string) int {
	opts := internal.GetOptStart(args)
	vars := internal.SetVars(opts)
	internal.Vars = vars.GetVars()

	g := internal.MakeProject{}
	g.SetSrc(opts.Tmpl)
	g.SetDest(opts.Project)
	g.SetTemp(opts.Tmpl)
	ret := g.Run()
	return ret
}
