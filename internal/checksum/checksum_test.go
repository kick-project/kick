package checksum

import (
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/utils"
)

func TestSha256sum(t *testing.T) {
	plainfile := filepath.Join(utils.TempDir(), "checksum", "plaintext.txt")
	origsumfile := filepath.Join(utils.TempDir(), "checksum", "plaintext.txt.sha256")
	newsumfile := filepath.Join(utils.TempDir(), "checksum", "plaintext.txt.sha256-test")
	Sha256sum(plainfile, newsumfile)

	passorig, sumorig, err := VerifySha256sum(plainfile, origsumfile)
	if err != nil {
		t.Error(err)
	}

	passnew, sumnew, err := VerifySha256sum(plainfile, newsumfile)
	if err != nil {
		t.Error(err)
	}

	if len(sumorig) == 0 || len(sumnew) == 0 {
		t.Fail()
	}
	if !passorig || !passnew {
		t.Fail()
	}
	if sumnew != sumorig {
		t.Fail()
	}
}

func TestVerifySha256sum(t *testing.T) {
	plainfile := filepath.Join(utils.TempDir(), "checksum", "plaintext.txt")
	sumfile := filepath.Join(utils.TempDir(), "checksum", "plaintext.txt.sha256")

	pass, sum, err := VerifySha256sum(plainfile, sumfile)
	if err != nil {
		t.Error(err)
	}
	if len(sum) == 0 {
		t.Fail()
	}
	if !pass {
		t.Fail()
	}
}
