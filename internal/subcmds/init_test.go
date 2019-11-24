package subcmds

import (
	"github.com/crosseyed/prjstart/internal"
	"log"
	"path/filepath"
)

func init() {
	home, _ := filepath.Abs("../../tmp/home")
	internal.Config = internal.LoadConfig(home, "")
	if internal.Config == nil {
		log.Fatalf("Config is nil\n\tHOME: %s\n", home)
	}
}
