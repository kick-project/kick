package utils

import "strings"

// SQL2fmt converts ? into %s and appends a newline to the string if it doesn't
// exist.
func SQL2fmt(sql string) (out string) {
	out = strings.ReplaceAll(sql, "?", "\"%s\"")
	out += "\n"
	return
}
