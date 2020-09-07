package compression

import (
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/checksum"
	"github.com/crosseyed/prjstart/internal/utils"
)

func TestCompress(t *testing.T) {
	src := filepath.Join(utils.TempDir(), "compression", "plaintext.txt")
	dst := filepath.Join(utils.TempDir(), "compression", "plaintext.txt.gz")
	sumfile := dst + ".sha256"
	sz, err := Compress(src, dst)
	if err != nil {
		t.Error(err)
	}
	pass, _, err := checksum.VerifySha256sum(dst, sumfile)
	switch {
	case sz == 0:
		t.Fail()
	case err != nil:
		t.Error(err)
		t.Fail()
	case !pass:
		t.Fail()
	}
}

func TestDecompress(t *testing.T) {
	var (
		sz  int64
		err error
	)
	src := filepath.Join(utils.TempDir(), "compression", "compressedtext.txt.gz")
	dst := filepath.Join(utils.TempDir(), "compression", "compressedtext.txt")
	sumfile := dst + ".sha256"
	sz, err = Decompress(src, dst)
	if err != nil {
		t.Error(err)
	}
	pass, _, err := checksum.VerifySha256sum(dst, sumfile)
	switch {
	case sz == 0:
		t.Fail()
	case err != nil:
		t.Error(err)
		t.Fail()
	case !pass:
		t.Fail()
	}
}
