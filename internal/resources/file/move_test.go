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
func TestMove_File(t *testing.T) {
	f, err := ioutil.TempFile("", "TestMove-*.txt")
	assert.Nil(t, err)
	_, _ = f.WriteString(`Original File`)
	f.Close()

	dest := filepath.Join(testtools.TempDir(), "TestMove-File-Target")
	err = file.Move(f.Name(), dest)
	assert.Nil(t, err)
	assert.FileExists(t, dest)
	assert.NoFileExists(t, f.Name())
	err = os.Remove(dest)
	if err != nil {
		t.Error(err)
	}
}

func TestMove_Dir(t *testing.T) {
	// Mock functions
	src, err := ioutil.TempDir("", "TestMove-Dir-Source-*")
	assert.Nil(t, err)

	dest := filepath.Join(testtools.TempDir(), "TestMove-Dir-Target")
	err = file.Move(src, dest)
	assert.Nil(t, err)
	assert.DirExists(t, dest)
	assert.NoFileExists(t, src)
	err = os.Remove(dest)
	if err != nil {
		t.Error(err)
	}
}

func TestMove_Recursive(t *testing.T) {
	// Mock functions
	src, err := ioutil.TempDir("", "TestMove-Recursive-*")
	assert.Nil(t, err)

	f1, err := ioutil.TempFile(src, "TestMove-File1-*")
	assert.Nil(t, err)
	f1.WriteString(`File1`) //nolint
	f1.Close()

	lvl1, err := ioutil.TempDir(src, "TestMove-Level1-*")
	assert.Nil(t, err)

	f2, err := ioutil.TempFile(lvl1, "TestMove-File2-*")
	assert.Nil(t, err)
	f2.WriteString(`File2`) //nolint
	f2.Close()

	dest := filepath.Join(testtools.TempDir(), "TestMove-Recursive-Target")
	err = file.Move(src, dest)
	assert.Nil(t, err)
	assert.DirExists(t, dest)
	assert.NoFileExists(t, src)
}
