package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//
// Tests
//
func TestGetTmpl(t *testing.T) {
	dwnlder := NewFetcher(Config)
	lp := dwnlder.GetTmpl("tmpl")
	assert.DirExists(t, lp)
}