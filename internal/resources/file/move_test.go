package file_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/resources/file"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/stretchr/testify/assert"
)

// TestMove test moving of files
func TestMove(t *testing.T) {
	f, err := ioutil.TempFile("", "TestMove-*.txt")
	assert.Nil(t, err)
	_, _ = f.WriteString(`Original File`)
	f.Close()

	dest := filepath.Join(testtools.TempDir(), "TestMove.txt")
	err = file.Move(f.Name(), dest)
	assert.Nil(t, err)
	assert.FileExists(t, dest)
	err = os.Remove(dest)
	if err != nil {
		t.Error(err)
	}
}
