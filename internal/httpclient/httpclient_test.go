package httpclient

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/checksum"
	"github.com/crosseyed/prjstart/internal/utils"
)

func TestDownload(t *testing.T) {
	var err error
	errchk := func() {
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test server
	servdir := filepath.Join(utils.TempDir(), "metadata", "serve")
	hdlr := http.FileServer(http.Dir(servdir))

	ts := httptest.NewServer(hdlr)
	defer ts.Close()

	gzfile := filepath.Join(utils.TempDir(), "metadata.db.gz")
	err = Download(ts.URL+"/metadata.db.gz", gzfile)
	errchk()

	shafile := filepath.Join(utils.TempDir(), "metadata.db.gz.sha256")
	err = Download(ts.URL+"/metadata.db.gz.sha256", shafile)
	errchk()

	pass, _, err := checksum.VerifySha256sum(gzfile, shafile)
	errchk()
	if !pass {
		t.Fail()
	}
}
