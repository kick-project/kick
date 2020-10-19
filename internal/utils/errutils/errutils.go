package errutils

import (
	"fmt"
	"log"

	"github.com/kick-project/kick/internal/utils"
)

func hasErrPrint(err error) bool {
	if err == nil {
		return false
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	e := log.Output(3, err.Error())
	if e != nil {
		panic(e)
	}
	log.SetFlags(log.LstdFlags)
	return true
}

func hasErrPrintf(format string, v ...interface{}) bool {
	hasError := false
	for _, e := range v {
		if _, ok := e.(error); ok {
			hasError = true
			break
		}
	}
	if !hasError {
		return false
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	out := fmt.Errorf(format, v...)
	log.Output(3, out.Error()) // nolint
	log.SetFlags(log.LstdFlags)
	return true
}

// Epanic will log an error and panic if err is not nil.
func Epanic(err error) {
	has := hasErrPrint(err)
	if !has {
		return
	}
	panic(err)
}

// Epanicf will log an error and panic if any argument passed to format is an error
func Epanicf(format string, v ...interface{}) {
	hasErr := hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	panic(fmt.Errorf(format, v...))
}

// Elogf will log an error if any argument passed to format is an error
func Elogf(format string, v ...interface{}) bool { // nolint
	return hasErrPrintf(format, v...)
}

// Efatal will log an error and exit if err is not nil.
func Efatal(err error) {
	has := hasErrPrint(err)
	if !has {
		return
	}
	utils.Exit(255)
}

// Efatalf will log an error and exit if any argument passed to fatal is an error
func Efatalf(format string, v ...interface{}) { // nolint
	hasErr := hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	utils.Exit(255)
}
