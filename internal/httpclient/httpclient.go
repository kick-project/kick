package httpclient

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/crosseyed/prjstart/internal/file"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

func checkUri(uri string) (valid bool, err error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return false, fmt.Errorf("Invalid URI %s: %w", uri, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false, fmt.Errorf("Unsupported URI scheme %s", u.Scheme)
	}
	return true, nil
}

func Download(uri string, fname string) error {
	ok, err := checkUri(uri)
	if !ok || err != nil {
		return err
	}
	resp, err := http.Get(uri)
	if errutils.Elogf(err, "Can not fetch URL %s: %v", uri, err) {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Can not fetch URL %s: Status code %d", uri, resp.StatusCode)
	}
	contentLength := resp.Header.Get("Content-Length")
	_ = contentLength

	file.Reader2File(resp.Body, fname)
	return nil
}
