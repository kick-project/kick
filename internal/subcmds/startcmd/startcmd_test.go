package startcmd_test

import (
	"testing"

	"github.com/kick-project/kick/internal/subcmds/startcmd"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", startcmd.UsageDoc)
}
