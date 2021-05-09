package errs

import (
	"fmt"

	"log"

	"github.com/kick-project/kick/internal/resources/exit"
)

// Errors error handling
type Errors struct {
	Ex     *exit.Handler `copier:"must"` // Exit handler
	Logger *log.Logger   `copier:"must"` // Default logger
}

// Panic will log an error and panic if err is not nil.
func (e *Errors) Panic(err error) {
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	panic(err)
}

// PanicF will log an error and panic if any argument passed to format is an error
func (e *Errors) PanicF(format string, v ...interface{}) {
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	panic(fmt.Errorf(format, v...))
}

// LogF will log an error if any argument passed to format is an error
func (e *Errors) LogF(format string, v ...interface{}) bool { // nolint
	return e.hasErrPrintf(format, v...)
}

// Fatal will log an error and exit if err is not nil.
func (e *Errors) Fatal(err error) {
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	e.Ex.Exit(255)
}

// FatalF will log an error and exit if any argument passed to fatal is an error
func (e *Errors) FatalF(format string, v ...interface{}) { // nolint
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	e.Ex.Exit(255)
}

func (e *Errors) hasErrPrint(err error) bool {
	if err == nil {
		return false
	}
	e.Logger.SetFlags(log.LstdFlags | log.Lshortfile)
	o := e.Logger.Output(3, err.Error())
	if o != nil {
		panic(o)
	}
	e.Logger.SetFlags(log.LstdFlags)
	return true
}

func (e *Errors) hasErrPrintf(format string, v ...interface{}) bool {
	hasError := false
	for _, elm := range v {
		if _, ok := elm.(error); ok {
			hasError = true
			break
		}
	}
	if !hasError {
		return false
	}
	e.Logger.SetFlags(log.LstdFlags | log.Lshortfile)
	out := fmt.Errorf(format, v...)
	e.Logger.Output(3, out.Error()) // nolint
	e.Logger.SetFlags(log.LstdFlags)
	return true
}