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

// TestCopy test coping of files
func TestCopy(t *testing.T) {
	f, err := ioutil.TempFile("", "TestCopy-*.txt")
	assert.Nil(t, err)
	_, _ = f.WriteString(`Original File`)
	f.Close()
	src := f.Name()

	dest := filepath.Join(testtools.TempDir(), "TestCopy.txt")
	_, err = file.Copy(src, dest)
	assert.Nil(t, err)

	for _, f := range []string{src, dest} {
		assert.FileExists(t, f)
		err = os.Remove(f)
		if err != nil {
			t.Error(err)
		}
	}
}
