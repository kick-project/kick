package subcmds

import (
	"log"
	"path/filepath"

	"github.com/crosseyed/prjstart/internal/config"
	"github.com/crosseyed/prjstart/internal/globals"
)

func init() {
	home, _ := filepath.Abs("../../tmp/home")
	globals.Config = config.Load(home, "")
	if globals.Config == nil {
		log.Fatalf("Config is nil\n\tHOME: %s\n", home)
	}
}
