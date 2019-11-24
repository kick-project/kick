package subcmds

import (
	"testing"
)

func TestList(t *testing.T) {
	args := []string{"list"}
	ret := List(args)
	_ = ret
}

func TestListLocal(t *testing.T) {
	t.Skip("Flag --local to be implemented")
	args := []string{"list", "--local"}
	ret := List(args)
	_ = ret
}

func TestListRemote(t *testing.T) {
	t.Skip("Flag --remote to be implemented")
	args := []string{"list", "--remote"}
	ret := List(args)
	_ = ret
}
