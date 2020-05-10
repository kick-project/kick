package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"sort"
)

func TestByVersion(t *testing.T) {
	target := []string{"3.0.0", "2.3.1", "2.0.1", "1.1.0", "0.1.1"}
	version := []string{"2.0.1", "0.1.1", "3.0.0", "1.1.0", "2.3.1"}
	bv := ByVersion(version)
	sort.Sort(sort.Reverse(bv))
	fmt.Println("*** Target ***")
	for _, ver := range target {
		fmt.Println(ver)
	}
	fmt.Println("\n*** Actual ***")
	for _, ver := range version {
		fmt.Println(ver)
	}

	assert.ElementsMatch(t, target, version)
}

func TestLatestVersion(t *testing.T) {
	target := "3.0.0"
	version := []string{"2.0.1", "0.1.1", "3.0.0", "1.1.0", "2.3.1"}
	actual := LatestVersion(version...)
	fmt.Printf("TARGET: %s\nACTUAL: %s\n", target, actual)
	assert.Equal(t, target, actual)
}