package subcmds

import "github.com/crosseyed/prjstart/internal"

func Install(args []string) int {
	opts := internal.GetOptInstall(args)
	_ = opts
	return 0
}
