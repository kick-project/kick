package errs

import (
	"fmt"
	"os"
	"sync"

	"log"

	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/logger"
)

// Handler error handling
//go:generate ifacemaker -f errs.go -s Handler -p errs -i HandlerIface -o errs_interfaces.go -c "AUTO GENERATED. DO NOT EDIT."
type Handler struct {
	ex     exit.HandlerIface  `validate:"required"` // Exit handler
	Logger logger.OutputIface `validate:"required"` // Default logger
	mu     *sync.Mutex
}

// New return a *Handler object.
func New(eh *exit.Handler, lgr logger.OutputIface) *Handler {
	return &Handler{
		ex:     eh,
		Logger: lgr,
		mu:     &sync.Mutex{},
	}
}

// Panic will log an error and panic if err is not nil.
func (e *Handler) Panic(err error) {
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	panic(err)
}

// PanicF will log an error and panic if any argument passed to format is an error
func (e *Handler) PanicF(format string, v ...interface{}) {
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	panic(fmt.Errorf(format, v...))
}

// LogF will log an error if any argument passed to format is an error
func (e *Handler) LogF(format string, v ...interface{}) bool { // nolint
	return e.hasErrPrintf(format, v...)
}

// Fatal will log an error and exit if err is not nil.
func (e *Handler) Fatal(err error) {
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	e.ex.Exit(255)
}

// FatalF will log an error and exit if any argument passed to fatal is an error
func (e *Handler) FatalF(format string, v ...interface{}) { // nolint
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	e.ex.Exit(255)
}

func (e *Handler) hasErrPrint(err error) bool {
	if err == nil {
		return false
	}
	o := e.Logger.Output(3, err.Error())
	if o != nil {
		panic(o)
	}
	return true
}

func (e *Handler) hasErrPrintf(format string, v ...interface{}) bool {
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
	out := fmt.Errorf(format, v...)
	e.Logger.Output(3, out.Error()) // nolint
	return true
}

// Panic will log an error and panic if err is not nil.
func Panic(err error) {
	e := makeErrors()
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	panic(err)
}

// PanicF will log an error and panic if any argument passed to format is an error
func PanicF(format string, v ...interface{}) {
	e := makeErrors()
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	panic(fmt.Errorf(format, v...))
}

// LogF will log an error if any argument passed to format is an error
func LogF(format string, v ...interface{}) bool { // nolint
	e := makeErrors()
	return e.hasErrPrintf(format, v...)
}

// Fatal will log an error and exit if err is not nil.
func Fatal(err error) {
	e := makeErrors()
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	e.ex.Exit(255)
}

// FatalF will log an error and exit if any argument passed to fatal is an error
func FatalF(format string, v ...interface{}) { // nolint
	e := makeErrors()
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	e.ex.Exit(255)
}

func makeErrors() *Handler {
	eh := &exit.Handler{
		Mode: exit.ExitMode,
	}
	e := &Handler{
		ex:     eh,
		Logger: logger.New(os.Stderr, "", log.LstdFlags, logger.ErrorLevel, eh),
	}
	return e
}
