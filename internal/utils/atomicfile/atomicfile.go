package atomicfile

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/crosseyed/prjstart/internal/utils"
)

func Reader2File(rdr io.Reader, dst string) (written int64, err error) {
	a := New(dst)
	written, err = a.Slurp(rdr)
	if err != nil {
		return written, err
	}
	err = a.Close()
	if err != nil {
		panic(err)
	}
	return written, err
}

type AtomicWrite struct {
	file    *os.File
	dst     string
	written int64
}

func New(dst string) *AtomicWrite {
	return &AtomicWrite{
		dst: dst,
	}
}

// Slurp Reads until EOF or an error occurs. Data is written to the tempfile
func (a *AtomicWrite) Slurp(rdr io.Reader) (written int64, err error) {
	f, err := a.tempfile()
	if err != nil {
		return 0, err
	}
	written, err = io.Copy(f, rdr)
	a.written += written
	return written, nil
}

// tempfile returns the *os.File object for the temporary file
func (a *AtomicWrite) tempfile() (*os.File, error) {
	if a.file != nil {
		return a.file, nil
	}
	f, err := ioutil.TempFile("", "")
	if utils.ChkErr(err, utils.Elogf, "Can not open temp file: %v", err) {
		return nil, err
	}
	a.file = f
	return a.file, nil
}

// Close closes the temporary file and moves to the destination
func (a *AtomicWrite) Close() error {
	if a.file == nil {
		return fmt.Errorf("Can not close file as file object is nil")
	}
	a.file.Close()
	os.Rename(a.file.Name(), a.dst)
	return nil
}
