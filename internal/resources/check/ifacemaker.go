// DO NOT EDIT: Generated using "make interfaces"

package check

// CheckIface ...
type CheckIface interface {
	// Init checks to see if an initialization has been performed. This function
	// will print an error message and exit if initialization is needed.
	Init() error
}
