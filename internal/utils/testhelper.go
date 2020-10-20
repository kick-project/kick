package utils

import (
	"path"
	"path/filepath"
	"runtime"
)

// FixtureDir returns the path to the fixture directory.
func FixtureDir() (fixturedir string) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Can not get filename")
	}
	fixturedir, err := filepath.Abs(path.Join(path.Dir(filename), "..", "..", "test", "fixtures"))
	if err != nil {
		panic(err)
	}
	return fixturedir
}

// TempDir returns the path to the localized tmp/ directory.
func TempDir() (tempdir string) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Can not get filename")
	}
	tempdir, err := filepath.Abs(path.Join(path.Dir(filename), "..", "..", "tmp"))
	if err != nil {
		panic(err)
	}
	return tempdir
}
