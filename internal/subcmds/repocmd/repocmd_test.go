package repocmd_test

import (
	"testing"

	"github.com/kick-project/kick/internal/subcmds/repocmd"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", repocmd.UsageDoc)
}
