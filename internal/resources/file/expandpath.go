package file

import (
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandPath expands the ~ into to the setting of the $HOME variable.
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = filepath.Join(dir, path[2:])
	}
	return path
}
