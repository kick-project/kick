// DO NOT EDIT: Generated using "make interfaces"

package errs

// HandlerIface ...
type HandlerIface interface {
	// Panic will log an error and panic if err is not nil.
	Panic(err error)
	// PanicF will log an error and panic if any argument passed to format is an error
	PanicF(format string, v ...interface{})
	// LogF will log an error if any argument passed to format is an error
	LogF(format string, v ...interface{}) bool
	// Fatal will log an error and exit if err is not nil.
	Fatal(err error)
	// FatalF will log an error and exit if any argument passed to fatal is an error
	FatalF(format string, v ...interface{})
}
