// AUTO GENERATED. DO NOT EDIT.

package start

// StartIface ...
type StartIface interface {
	// Start start command
	Start(projectname, template, path string)
	// List lists the output
	List(long bool)
	// Show show files used in a template. base is the path to the template
	// directory on local disk. If a slice of incLabels is provided, only files,
	// directories that have matching labels will be displayed. If a file or
	// directory has no label it is always displayed.
	Show(hdle string, incLabels []string, ops ShowOptions)
}
