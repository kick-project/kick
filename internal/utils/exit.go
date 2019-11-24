package utils

import (
	"fmt"
	"os"
)

const (
	Tnone = iota
	Ttesting
)

var Mode int

func Exit(code int) {
	switch Mode {
	case Tnone:
		os.Exit(code)
	case Ttesting:
		panic(fmt.Sprintf("Exit %d\n", code))
	default:
		panic(fmt.Sprintf("Unknown testing mode: %d", Mode))
	}
}
