package compression

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/kick-project/kick/internal/utils/errutils"
)

// Compress takes the path plainSrc compresses it and writes it to the path gzDst.
// Returns the number of bytes that were written to gzDst.
func Compress(plainSrc, gzDst string) (written int64, err error) {
	srcIO, err := os.Open(plainSrc)
	if err != nil {
		return 0, fmt.Errorf("Failed to open file %s: %s", plainSrc, err)
	}

	dstIO, err := ioutil.TempFile("", "")
	if err != nil {
		return 0, fmt.Errorf("Failed to open temporary file for %s: %s", gzDst, err)
	}

	gzWriter := gzip.NewWriter(dstIO)

	// TODO: Bugfix written should reflect the number of bytes written not read.
	written, err = io.Copy(gzWriter, srcIO)
	if err != nil {
		return 0, err
	}
	srcIO.Close()
	gzWriter.Close()
	dstIO.Close()

	err = os.Rename(dstIO.Name(), gzDst)
	errutils.Epanic(err)
	return written, nil
}

// Decompress loads the path gzSrc decompresses the contents and writes to plainDst.
// Returns the number of bytes that were written to plainDst.
func Decompress(gzSrc, plainDst string) (written int64, err error) {
	srcIO, err := os.Open(gzSrc)
	if err != nil {
		return 0, fmt.Errorf("Failed to open file %s: %s", gzSrc, err)
	}
	dstIO, err := ioutil.TempFile("", "")
	if err != nil {
		return 0, fmt.Errorf("Failed to open temporary file for %s: %s", plainDst, err)
	}

	gzReader, err := gzip.NewReader(srcIO)
	if err != nil {
		return 0, err
	}

	// TODO: Bugfix written should reflect the number of bytes written not read.
	written, err = io.Copy(dstIO, gzReader)
	if err != nil {
		return 0, err
	}
	srcIO.Close()
	err = gzReader.Close()
	errutils.Epanic(err)
	dstIO.Close()

	err = os.Rename(dstIO.Name(), plainDst)
	errutils.Epanic(err)
	return written, nil
}