package errutils

import (
	"fmt"
	"log"

	"github.com/crosseyed/prjstart/internal/utils"
)

func logErr(format string, v ...interface{}) bool {
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

// Epanicf will log an error and panic if any argument passed to format is an error
func Epanicf(format string, v ...interface{}) {
	hasErr := logErr(format, v...)
	if !hasErr {
		return
	}
	panic(fmt.Errorf(format, v...))
}

// Elogf will log an error if any argument passed to format is an error
func Elogf(format string, v ...interface{}) bool { // nolint
	return logErr(format, v...)
}

// Efatalf will log an error and exit if any argument passed to fatal is an error
func Efatalf(format string, v ...interface{}) { // nolint
	hasErr := logErr(format, v...)
	if !hasErr {
		return
	}
	utils.Exit(255)
}
