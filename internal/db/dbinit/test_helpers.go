package dbinit

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func InitTestDB() string {
	_, b, _, _ := runtime.Caller(0)
	p, err := filepath.Abs(filepath.Join(filepath.Dir(b), "..", "..", "..", "tmp", "Test.db"))
	if err != nil {
		panic(fmt.Sprintf("filepath.Abs %s error: %v", p, err))
	}
	os.Remove(p)
	i := New(p)
	i.Init()
	return p
}
