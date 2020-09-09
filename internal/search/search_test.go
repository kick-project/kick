package search

import (
	"bytes"
	"testing"
)

func TestSearch(t *testing.T) {
	b := bytes.Buffer{}
	s := Search{}
	s.Search("template1", &b)
	strout := b.String()
	if len(strout) == 0 {
		t.Fail()
	}
}
