package initcmd_test

import (
	"testing"

	"github.com/kick-project/kick/internal/subcmds/initcmd"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", initcmd.UsageDoc)
}
