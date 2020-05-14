package utils

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByVersion(t *testing.T) {
	target := []string{"3.0.0", "2.3.1", "2.0.1", "1.1.0", "0.1.1"}
	version := []string{"2.0.1", "0.1.1", "3.0.0", "1.1.0", "2.3.1"}
	bv := ByVersion(version)
	sort.Sort(sort.Reverse(bv))
	assert.ElementsMatch(t, target, version)
}

func TestLatestVersion(t *testing.T) {
	target := "3.0.0"
	version := []string{"2.0.1", "0.1.1", "3.0.0", "1.1.0", "2.3.1"}
	actual := LatestVersion(version...)
	assert.Equal(t, target, actual)
}
