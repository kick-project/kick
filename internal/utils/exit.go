package utils

import (
	"fmt"
	"os"
)

const (
	// MNone set the exit mode at default. When utils.Exit is called it exits
	// with the as normal with the supplied return code.
	MNone = iota
	// MPanic set the exit mode to panic. When utils.Exit is called then a panic
	// will ensue instead of exiting.
	MPanic
)

var exitMode int

// ExitMode sets the exit mode for the utils.Exit call. Current modes are
//
// * utils.MNone (default) - call exit mode directly
//
// * utils.MPanic - Panic instead of exiting
//
// Returns true if exit mode was was successfully set.
func ExitMode(mode int) (ok bool) {
	switch mode {
	case MNone:
		exitMode = MNone
		ok = true
	case MPanic:
		exitMode = MPanic
		ok = true
	}
	return
}

// Exit will either exit (default) or panic. This behaviour is controlled by
// ExitMode.
func Exit(code int) {
	switch exitMode {
	case MNone:
		os.Exit(code)
	case MPanic:
		panic(fmt.Sprintf("Exit %d\n", code))
	default:
		panic(fmt.Sprintf("Unknown testing mode: %d", exitMode))
	}
}
