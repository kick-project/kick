package utils

import (
	"fmt"
	"log"
)

func ChkErr(e error, f func(e error, v ...interface{}), v ...interface{}) bool {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e != nil {
		f(e, v...)
		return true
	}
	log.SetFlags(log.LstdFlags)
	return false
}

func Epanicf(e error, args ...interface{}) {
	if e == nil {
		return
	}
	var format string
	var out string
	if len(args) > 0 {
		f, ok := args[0].(string)
		args = args[:len(args)-1]
		if ok {
			format = f
		}
		out = fmt.Sprintf(format, args...)
		log.Output(2, out) // nolint
	}
	panic(out)
}

func Elogf(e error, args ...interface{}) { // nolint
	if e == nil {
		return
	}
	var format string
	var out string
	if len(args) > 0 {
		if f, ok := args[0].(string); ok {
			format = f
			args = args[1:]
		}
		out = fmt.Sprintf(format, args...)
		log.Output(3, out) // nolint
	}
}

func Efatalf(e error, args ...interface{}) { // nolint
	if e == nil {
		return
	}
	var format string
	var out string
	if len(args) > 0 {
		if f, ok := args[0].(string); ok {
			format = f
			args = args[1:]
		}
		out = fmt.Sprintf(format, args...)
		log.Output(2, out) // nolint
	}
	Exit(255)
}
