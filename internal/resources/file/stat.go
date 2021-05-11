package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/kick-project/kick/internal/resources/errs"
)

// A Info describes a file and is returned by Stat.
type Info interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() fs.FileMode  // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
	Path() string       // path to file
	Abs() string        // absolute path to file
	Exists() bool       // true if file exists
}

// info file info
type info struct {
	stat    fs.FileInfo
	path    string
	statErr error
}

// Name basename of the file
func (i *info) Name() string {
	if i.stat != nil {
		return i.stat.Name()
	}
	return ""
}

// Size length of file in bytes for regular files, system dependent for others. Returns -1 on error
func (i *info) Size() int64 {
	if i.stat != nil {
		return i.stat.Size()
	}
	return -1
}

// Mode file mode bits, Returns 0 on error
func (i *info) Mode() fs.FileMode {
	if i.stat != nil {
		return i.stat.Mode()
	}
	return 0

}

// ModTime modification time. Returns an unpopulated struct on error
func (i *info) ModTime() time.Time {
	if i.stat != nil {
		return i.stat.ModTime()
	}
	return time.Time{}
}

// IsDir abbreviation for Mode().IsDir()
func (i *info) IsDir() bool {
	if i.stat != nil {
		return i.stat.IsDir()
	}
	return false
}

// Sys underlying data source (can return nil)
func (i *info) Sys() interface{} {
	if i.stat != nil {
		return i.stat.Sys()
	}
	return nil
}

// Path path to file
func (i *info) Path() string {
	return i.path
}

// Abs absolute path
func (i *info) Abs() string {
	p, err := filepath.Abs(i.path)
	errs.Panic(err)
	return p
}

// Exists true if file exists
func (i *info) Exists() bool {
	if i.stat != nil {
		return true
	} else if i.statErr != nil && os.IsNotExist(i.statErr) {
		return false
	} else {
		errs.PanicF(`error stating file %s: %v`, i.path, i.statErr)
	}
	return false
}

// Stat wrapper around stat
func Stat(path string) (Info, error) {
	i := &info{
		path: path,
	}
	i.stat, i.statErr = os.Stat(path)

	return i, i.statErr
}
