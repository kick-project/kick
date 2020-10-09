package strs

import (
	"bufio"
	"strings"
)

// SplitLines split a string into lines
func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
