package subcmds

import (
	"path/filepath"

	"github.com/crosseyed/prjstart/internal"
	"github.com/crosseyed/prjstart/internal/services/template"
	"github.com/crosseyed/prjstart/internal/settings"
)

// Start start cli option
func Start(args []string, s *settings.Settings) int {
	opts := internal.GetOptStart(args)
	varsGen := internal.SetVars(opts)
	vars := varsGen.GetVars()
	_ = vars

	tOpts := *s.Template()

	// Set project name
	name := filepath.Base(opts.Project)
	tOpts.Variables.SetProjectVar("NAME", name)
	t := template.New(tOpts)
	t.SetSrcDest(opts.Tmpl, opts.Project)
	ret := t.Run()
	return ret
}
