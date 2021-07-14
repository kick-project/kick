// AUTO GENERATED. DO NOT EDIT.

package initialize

// InitIface ...
type InitIface interface {
	// CreateRepo create repository
	CreateRepo(name, path string) int
	// CreateTemplate create template
	CreateTemplate(name, path string) int
}
