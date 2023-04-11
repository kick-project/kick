package file

import (
	"bytes"
	_ "embed"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed readfile.txt
var readFileText []byte
var expected = readFileText

var readFile = "readfile.txt"

func TestReadFile_Byte(t *testing.T) {
	actual, err := ReadFile(readFile, readFileText)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestReadFile_String(t *testing.T) {
	actual, err := ReadFile(readFile, string(readFileText))
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestReadFile_BytesBuffer(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write(readFileText)
	actual, err := ReadFile(readFile, buf)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestReadFile_IoWriter(t *testing.T) {
	rdr := bytes.NewReader(readFileText)
	actual, err := ReadFile(readFile, io.Reader(rdr))
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestReadFile_File(t *testing.T) {
	actual, err := ReadFile(readFile, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
