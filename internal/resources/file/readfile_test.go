package file

import (
	"bytes"
	_ "embed"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed readfile.txt
var readFileText []byte
var expected = readFileText

var readFile = "readfile.txt"

func TestReadLines_Byte_unlimited(t *testing.T) {
	actual, err := ReadLines(readFile, readFileText, 10)
	assert.NoError(t, err)
	assert.Equal(t, 5, strings.Count(string(actual), "\n"))
}

func TestReadLines_String_unlimited(t *testing.T) {
	actual, err := ReadLines(readFile, string(readFileText), 10)
	assert.NoError(t, err)
	assert.Equal(t, 5, strings.Count(string(actual), "\n"))
}

func TestReadLines_BytesBuffer_unlimited(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write(readFileText)
	actual, err := ReadLines(readFile, buf, 10)
	assert.NoError(t, err)
	assert.Equal(t, 5, strings.Count(string(actual), "\n"))
}

func TestReadLines_IoReader_unlimited(t *testing.T) {
	rdr := bytes.NewReader(readFileText)
	actual, err := ReadLines(readFile, io.Reader(rdr), 10)
	assert.NoError(t, err)
	assert.Equal(t, 5, strings.Count(string(actual), "\n"))
}

func TestReadLines_Byte(t *testing.T) {
	actual, err := ReadLines(readFile, readFileText, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, strings.Count(string(actual), "\n"))
}

func TestReadLines_String(t *testing.T) {
	actual, err := ReadLines(readFile, string(readFileText), 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, strings.Count(string(actual), "\n"))
}

func TestReadLines_BytesBuffer(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write(readFileText)
	actual, err := ReadLines(readFile, buf, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, strings.Count(string(actual), "\n"))
}

func TestReadLines_IoReader(t *testing.T) {
	rdr := bytes.NewReader(readFileText)
	actual, err := ReadLines(readFile, io.Reader(rdr), 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, strings.Count(string(actual), "\n"))
}

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

func TestReadFile_IoReader(t *testing.T) {
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
