package internal

import (
	"path"
	"path/filepath"
	"runtime"
)

func FixtureDir() (fixturedir string) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("Can not get filename")
	}
	fixturedir, err := filepath.Abs(path.Join(path.Dir(filename), "..", "..", "test", "fixtures"))
	if err != nil {
		panic(err)
	}
	return fixturedir
}
