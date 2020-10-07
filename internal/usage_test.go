package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOptMainStart(t *testing.T) {
	args := []string{"prjstart", "start", "template", "project"}
	o := GetOptMain(args)
	assert.True(t, o.Start)
	assert.False(t, o.Install)
	assert.False(t, o.List)
}

func TestGetOptMainList(t *testing.T) {
	args := []string{"prjstart", "list"}
	o := GetOptMain(args)
	assert.True(t, o.List)
	assert.False(t, o.Start)
	assert.False(t, o.Install)
}

func TestGetOptMainInstall(t *testing.T) {
	t.Skip("install - to be implemented")
	args := []string{"prjstart", "install"}
	o := GetOptMain(args)
	t.Skip()
	assert.True(t, o.Install)
	assert.False(t, o.Start)
	assert.False(t, o.List)
}
