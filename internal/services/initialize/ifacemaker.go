// DO NOT EDIT: Generated using "make interfaces"

package initialize

// InitIface ...
type InitIface interface {
	// CreateRepo create repository
	CreateRepo(name string) int
	// CreateTemplate create template
	CreateTemplate(name string) int
}
