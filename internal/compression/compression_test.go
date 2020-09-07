package compression

import (
	"fmt"
	"testing"

	"github.com/crosseyed/prjstart/internal/db/dbinit"
)

func TestCompress(t *testing.T) {
	src := dbinit.InitTestDB()
	dst := fmt.Sprintf("%s.gz", src)
	sz, err := Compress(src, dst)
	switch {
	case sz == 0:
		t.Fail()
	case err != nil:
		t.Error(err)
		t.Fail()
	}
}

func TestDecompress(t *testing.T) {
	var (
		sz  int64
		err error
	)
	chkErr := func() {
		switch {
		case sz == 0:
			t.Fail()
		case err != nil:
			t.Error(err)
			t.Fail()
		}
	}
	src := dbinit.InitTestDB()
	dst := fmt.Sprintf("%s.gz", src)
	sz, err = Compress(src, dst)
	chkErr()
	sz, err = Decompress(dst, src)
	chkErr()
}
