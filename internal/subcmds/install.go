package subcmds

import (
	"github.com/crosseyed/prjstart/internal"
	"github.com/crosseyed/prjstart/internal/settings"
)

func Install(args []string, s *settings.Settings) int {
	opts := internal.GetOptInstall(args)
	_ = opts
	_ = s
	return 0
}
