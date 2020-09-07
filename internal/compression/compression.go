package compression

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func Compress(plainSrc, gzDst string) (bytes int64, err error) {
	srcIO, err := os.Open(plainSrc)
	if err != nil {
		return 0, fmt.Errorf("Failed to open file %s: %s", plainSrc, err)
	}

	dstIO, err := ioutil.TempFile("", "")
	if err != nil {
		return 0, fmt.Errorf("Failed to open temporary file for %s: %s", gzDst, err)
	}

	gzWriter := gzip.NewWriter(dstIO)

	written, err := io.Copy(gzWriter, srcIO)
	if err != nil {
		return 0, err
	}
	srcIO.Close()
	gzWriter.Close()
	dstIO.Close()

	os.Rename(dstIO.Name(), gzDst)
	return written, nil
}

func Decompress(gzSrc, plainDst string) (bytes int64, err error) {
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
	written, err := io.Copy(dstIO, gzReader)
	if err != nil {
		return 0, err
	}
	srcIO.Close()
	gzReader.Close()
	dstIO.Close()

	os.Rename(dstIO.Name(), plainDst)
	return written, nil
}
