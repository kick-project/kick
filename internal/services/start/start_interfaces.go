// AUTO GENERATED. DO NOT EDIT.

package start

// StartIface ...
type StartIface interface {
	// Start start command
	Start(projectname, template, path string)
	// List lists the output
	List(long bool)
	// Show show files
	Show(base string, filter []string)
}
