package errutils

import (
	"fmt"
	"log"

	"github.com/crosseyed/prjstart/internal/utils"
)

func logErr(err error, format string, v ...interface{}) bool {
	if err == nil {
		return false
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	out := fmt.Sprintf(format, v...)
	log.Output(3, out) // nolint
	log.SetFlags(log.LstdFlags)
	return true
}

func Epanicf(err error, format string, v ...interface{}) bool {
	hasErr := logErr(err, format, v...)
	if !hasErr {
		return false
	}
	panic(err)
	return true
}

func Elogf(err error, format string, v ...interface{}) bool { // nolint
	return logErr(err, format, v...)
}

func Efatalf(err error, format string, v ...interface{}) bool { // nolint
	hasErr := logErr(err, format, v...)
	if !hasErr {
		return hasErr
	}
	utils.Exit(255)
	return hasErr
}
