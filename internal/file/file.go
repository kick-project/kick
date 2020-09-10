package file

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

func Reader2File(rdr io.Reader, dst string) (written int64, err error) {
	a := NewAtomicWrite(dst)
	written, err = a.Copy(rdr)
	if err != nil {
		return written, err
	}
	err = a.Close()
	if err != nil {
		panic(err)
	}
	return written, err
}

// AtomicWrite atomically writes files by using a temp file.
// When Close is called the temp file is closed and moved to its final destination.
type AtomicWrite struct {
	file    *os.File
	dst     string
	written int64
}

func NewAtomicWrite(dst string) *AtomicWrite {
	return &AtomicWrite{
		dst: dst,
	}
}

// Close closes the temporary file and moves to the destination
func (a *AtomicWrite) Close() error {
	if a.file == nil {
		err := fmt.Errorf("Object is nil")
		errutils.Elogf("Can not close file: %w", err)
	}
	a.file.Close()
	os.Rename(a.file.Name(), a.dst)
	return nil
}

// Copy Reads until EOF or an error occurs. Data is written to the tempfile
func (a *AtomicWrite) Copy(rdr io.Reader) (written int64, err error) {
	f, err := a.tempfile()
	if err != nil {
		return 0, err
	}
	written, err = io.Copy(f, rdr)
	a.written += written
	return written, nil
}

// Write writes bytes to the tempfile
func (a *AtomicWrite) Write(data []byte) (written int, err error) {
	f, err := a.tempfile()
	if err != nil {
		return 0, err
	}
	written, err = f.Write(data)
	if errutils.Elogf("Can not write to temporary file: %w", err) {
		return written, err
	}
	return written, nil
}

// tempfile returns the *os.File object for the temporary file
func (a *AtomicWrite) tempfile() (*os.File, error) {
	if a.file != nil {
		return a.file, nil
	}
	f, err := ioutil.TempFile("", "")
	if errutils.Elogf("Can not open temp file: %v", err) {
		return nil, err
	}
	a.file = f
	return a.file, nil
}
