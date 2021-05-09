package checksum_test

import (
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/resources/checksum"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/utils"
)

func TestSha256sum(t *testing.T) {
	plainfile := filepath.Join(utils.TempDir(), "checksum", "plaintext.txt")
	origsumfile := filepath.Join(utils.TempDir(), "checksum", "plaintext.txt.sha256")
	newsumfile := filepath.Join(utils.TempDir(), "checksum", "plaintext.txt.sha256-test")
	_, err := checksum.Sha256SumFile(plainfile, newsumfile)
	errs.Panic(err)

	passorig, sumorig, err := checksum.VerifySha256sum(plainfile, origsumfile)
	if err != nil {
		t.Error(err)
	}

	passnew, sumnew, err := checksum.VerifySha256sum(plainfile, newsumfile)
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

	pass, sum, err := checksum.VerifySha256sum(plainfile, sumfile)
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
