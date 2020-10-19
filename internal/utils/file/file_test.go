package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kick-project/kick/internal/utils"
)

func TestReader2File(t *testing.T) {
	rdr := strings.NewReader("Test data\n")
	dst := filepath.Join(utils.TempDir(), "TestReader2File.txt")
	written, err := Reader2File(rdr, dst)
	if err != nil {
		t.Error(err)
	}
	if written == 0 {
		t.Fail()
	}
	if info, err := os.Stat(dst); os.IsNotExist(err) || info.Size() == 0 {
		t.Fail()
	}
}
