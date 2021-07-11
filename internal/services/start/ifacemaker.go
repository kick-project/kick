// DO NOT EDIT: Generated using "make interfaces"

package start

// StartIface ...
type StartIface interface {
	// Start start command
	Start(projectname, template, path string)
	// List lists the output
	List(long bool) int
}
