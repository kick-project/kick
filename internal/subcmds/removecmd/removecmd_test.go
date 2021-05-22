package removecmd_test

import (
	"testing"

	"github.com/kick-project/kick/internal/subcmds/removecmd"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", removecmd.UsageDoc)
}
