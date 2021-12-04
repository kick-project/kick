package atomicfile_test

import (
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kick-project/kick/internal/resources/file/atomicfile"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/stretchr/testify/assert"
)

// TestClose
func TestAtomicFile_Close_error(t *testing.T) {
	// Path to target file
	dstfile := filepath.Join(testtools.TempDir(), "atomicfile_close.txt")

	f := atomicfile.New(dstfile)
	err := f.Close()

	assert.Error(t, err)
}

// TestCopy
func TestAtomicFile_Copy(t *testing.T) {
	contentlen := 1024
	contentsrc := make([]byte, contentlen)
	contentdst := make([]byte, contentlen)
	dstfile := filepath.Join(testtools.TempDir(), "atomicfile_copy.txt")

	// Save data to tempfile
	tmpfile, _ := ioutil.TempFile("", "atomicfile_copy*.txt") // nolint
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}()

	// Random string
	rand.Seed(time.Now().UnixNano())
	letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := range contentsrc {
		contentsrc[i] = letters[rand.Intn(len(letters))]
	}
	tmpfile.Write(contentsrc) // nolint

	// Setup temp file for reading
	tmpfile.Sync()                // nolint
	tmpfile.Seek(0, io.SeekStart) // nolint

	// New AtomicFile
	f := atomicfile.New(dstfile)
	written, err := f.Copy(tmpfile)
	if err != nil {
		t.Error(err)
	}
	f.Close() // nolint

	// Open destination file for reading
	readfile, _ := os.Open(dstfile) // nolint
	readfile.Read(contentdst)       // nolint

	assert.Equal(t, written, int64(contentlen))
	assert.Equal(t, contentsrc, contentdst)
}
