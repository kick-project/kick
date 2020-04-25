package subcmds

import (
	"testing"
)

func TestList(t *testing.T) {
	args := []string{"list"}
	ret := List(args)
	_ = ret
}
