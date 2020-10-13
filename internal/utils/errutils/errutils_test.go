package errutils

import (
	"errors"
	"testing"
)

func TestHasErrPrintf(t *testing.T) {
	if !hasErrPrintf("Error Test: %w", errors.New("This is an error")) {
		t.Fail()
	}
}

func TestHasErrPrint(t *testing.T) {
	if !hasErrPrintf("Error Test: %w", errors.New("This is an error")) {
		t.Fail()
	}
}
