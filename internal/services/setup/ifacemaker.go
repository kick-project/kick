// DO NOT EDIT: Generated using "make interfaces"

package setup

// SetupIface ...
type SetupIface interface {
	// Init initialize everything.
	Init()
	// InitPaths initialize paths.
	InitPaths()
	// InitMetadata initialize metadata.
	InitMetadata()
	// InitConfig initialize configuration file.
	InitConfig()
}
