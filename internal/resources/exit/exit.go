// Package exit handles exit to OS
package exit

import (
	"fmt"
	"os"
)

const (
	// MNone set the exit mode at default. When Exit.Exit is called it exits
	// with the as normal with the supplied return code.
	MNone = iota
	// MPanic set the exit mode to panic. When Exit.Exit is called then a panic
	// will ensue instead of exiting.
	MPanic
)

var (
	exitMode int
)

// Handler exit handling
type Handler struct {
	Mode int
}

// Exit exit function with exit code. 0 is success
func (e *Handler) Exit(code int) {
	switch e.Mode {
	case MNone:
		os.Exit(code)
	case MPanic:
		panic(fmt.Sprintf("Exit %d\n", code))
	default:
		panic(fmt.Sprintf("Unknown exit mode: %d", e.Mode))
	}
}

// Exit exit function with exit code. 0 is success
func Exit(code int) {
	h := Handler{
		Mode: exitMode,
	}
	h.Exit(code)
}

// Mode sets the exit mode for the exit.Exit call. Current modes are,
// exit.None (default - call exit mode directly) & exit.Panic (panic instead of exiting)
// Returns true if exit mode was was successfully set.
func Mode(mode int) (ok bool) {
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
