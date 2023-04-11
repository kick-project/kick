package modeline_test

import (
	"testing"

	"github.com/kick-project/kick/internal/resources/modeline"
	"github.com/stretchr/testify/assert"
)

var parseFile = `file.txt`
var parseTest = `# kick:render ignore label=core,editor`

func TestParser(t *testing.T) {
	ml, err := modeline.Parse(parseFile, parseTest, 1)
	assert.NoError(t, err)
	assert.True(t, ml.Option("render"))
	assert.True(t, ml.Option("ignore"))
	assert.True(t, ml.Label("core"))
	assert.True(t, ml.Label("editor"))
}
