package dbinit

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitTestDB() string {
	p, err := filepath.Abs(filepath.Join("..", "..", "..", "tmp", "Test.db"))
	if err != nil {
		panic(fmt.Sprintf("filepath.Abs %s error: %v", p, err))
	}
	os.Remove(p)
	i := New(p)
	i.Init()
	return p
}
