// DO NOT EDIT: Generated using "make interfaces"

package update

import (
	_ "github.com/mattn/go-sqlite3"
)

// UpdateIface ...
type UpdateIface interface {
	// Build metadata. Conf defaults to globals.Config if Conf is nil.
	Build() error
}
