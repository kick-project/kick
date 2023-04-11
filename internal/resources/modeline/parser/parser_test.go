package parser

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var parseFile = `file.txt`
var parseTest = `# kick:render type=core,editor`
var parseTestNoML = "# 1\n# 2\n# 3"
var parseTestOutOfRange = "#\n#\n#\n# kick:render"

func TestParse(t *testing.T) {
	items := Parse(parseFile, parseTest, 5)
	assert.Contains(t, items, Item{Type: OPTION, Value: "render"})
	assert.Contains(t, items, Item{Type: OPTION, Value: "type"})
	assert.Contains(t, items, Item{Type: TYPE, Value: "core"})
	assert.Contains(t, items, Item{Type: TYPE, Value: "editor"})
}

func TestParse_NoModeLine(t *testing.T) {
	items := Parse(parseFile, parseTestNoML, 5)
	assert.Empty(t, items)
}

func TestParse_OutOfBounds(t *testing.T) {
	items := Parse(parseFile, parseTestOutOfRange, 3)
	assert.Empty(t, items)
}
