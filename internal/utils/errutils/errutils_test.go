package errutils

import (
	"errors"
	"testing"
)

func TestLogErr(t *testing.T) {
	if !logErr("Error Test: %w", errors.New("This is an error")) {
		t.Fail()
	}
}
