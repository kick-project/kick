package utils

import (
	"fmt"
	"log"
)

func ChkErr(e error, f func(error, string, ...interface{}), format string, v ...interface{}) bool {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e != nil {
		f(e, format, v...)
		return true
	}
	log.SetFlags(log.LstdFlags)
	return false
}

func Epanicf(e error, format string, args ...interface{}) {
	if e == nil {
		return
	}
	out := fmt.Sprintf(format, args...)
	log.Output(3, out) // nolint
	panic(out)
}

func Elogf(e error, format string, args ...interface{}) { // nolint
	out := fmt.Sprintf(format, args...)
	log.Output(3, out) // nolint
}

func Efatalf(e error, format string, args ...interface{}) { // nolint
	if e == nil {
		return
	}
	out := fmt.Sprintf(format, args...)
	log.Output(3, out) // nolint
	Exit(255)
}
