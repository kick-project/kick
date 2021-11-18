package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type MockMoveFunc func(src, dst string) error

var (
	MockMove MockMoveFunc
)

// Move recursively move a directory or file from src to dst.
func Move(src, dst string) error {
	if MockMove != nil {
		return MockMove(src, dst)
	}
	return MoveAll(src, dst)
}

// MoveAll recursively move a directory or file from src to dst.
func MoveAll(src, dst string) error {
	err := CopyAll(src, dst)
	if err != nil {
		return err
	}
	err = os.RemoveAll(src)
	return err

}

// CopyAll recursively copy a directory or file from src to dst.
func CopyAll(src, dst string) error {
	/* Source Checks */

	// Check src exists
	srcInfo, err := os.Stat(src)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("move %s %s: no such file or directory", src, dst)
	} else if err != nil {
		return err
	}

	// Check file is a regular file
	if !srcInfo.Mode().IsRegular() && !srcInfo.Mode().IsDir() {
		return fmt.Errorf("move %s %s: not a file or directory", src, dst)
	}

	// Check parent directory of destination exits
	dir := filepath.Dir(dst)
	dstInfo, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("move %s %s: no such file or directory", src, dst)
	}

	// Check if parent of destination is a directory
	if dstInfo.Mode()&os.ModeDir != os.ModeDir {
		return fmt.Errorf("move %s %s: parent of destination is not a directory", src, dst)
	}

	err = filepath.Walk(src, func(srcPath string, srcPathInfo os.FileInfo, err error) error {
		return walkFunc(src, dst, srcPath, srcPathInfo, err)
	})

	return err
}

func walkFunc(src, dst, srcPath string, srcPathInfo os.FileInfo, err error) error {
	srcMode := srcPathInfo.Mode()
	dstPath := filepath.Join(dst, strings.TrimPrefix(srcPath, src))
	switch {
	case errors.Is(err, os.ErrNotExist):
		return err
	case srcMode.IsDir():
		err = os.MkdirAll(dstPath, os.ModePerm)
		return err
	case srcMode.IsRegular():
		_, err := CopyFile(srcPath, dstPath)
		return err
	case srcMode&os.ModeSymlink != 0:
		link, err := os.Readlink(srcPathInfo.Name())
		if err != nil {
			return err
		}
		err = os.Symlink(dstPath, link)
		return err
	}
	return fmt.Errorf("move %s %s: unsupported file type", srcPath, dstPath)
}

// MoveFile move src file to dst, returns the number of bytes that were moved.
func MoveFile(src, dst string) (int64, error) {
	n, err := Copy(src, dst)
	if err != nil {
		return n, err
	}
	err = os.Remove(src)
	return n, err
}

// CopyFile copy src file to dst, returns the number of bytes that were copied.
func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
