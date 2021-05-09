package compression_test

import (
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/resources/checksum"
	"github.com/kick-project/kick/internal/resources/compression"
	"github.com/kick-project/kick/internal/resources/testtools"
)

func TestCompress(t *testing.T) {
	src := filepath.Join(testtools.TempDir(), "compression", "plaintext.txt")
	dst := filepath.Join(testtools.TempDir(), "compression", "plaintext.txt.gz")
	sumfile := dst + ".sha256"
	sz, err := compression.Compress(src, dst)
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
	src := filepath.Join(testtools.TempDir(), "compression", "compressedtext.txt.gz")
	dst := filepath.Join(testtools.TempDir(), "compression", "compressedtext.txt")
	sumfile := dst + ".sha256"
	sz, err = compression.Decompress(src, dst)
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
