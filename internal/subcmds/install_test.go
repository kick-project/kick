package subcmds

import (
	"testing"
)

func TestInstall(t *testing.T) {
	args := []string{"install", "set1", "tmpl1"}
	Install(args)
}
