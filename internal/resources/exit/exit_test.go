package exit_test

import (
	"fmt"
	"testing"

	"github.com/kick-project/kick/internal/resources/exit"
)

func Test_Exit(t *testing.T) {
	exit.Mode(exit.Panic)

	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected a panic")
		}
	}()
	exit.Exit(255)
}

func TestExit_Exit_Panic(t *testing.T) {
	e := exit.Handler{
		Mode: exit.Panic,
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected a panic")
		}
	}()
	e.Exit(255)
}

func TestExit_Exit_Unknown(t *testing.T) {
	e := exit.Handler{
		Mode: 5,
	}

	msg := fmt.Sprintf("Unknown exit mode: %d", e.Mode)
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected a panic")
		} else if r.(string) != msg {
			t.Fail()
		}
	}()
	e.Exit(255)
}
