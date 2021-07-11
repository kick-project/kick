package checkvars_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/stretchr/testify/assert"
)

func TestCheck_Check(t *testing.T) {
	in := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	inject := di.New(&di.Options{
		Home:   filepath.Join(testtools.TempDir(), "home"),
		Stdout: stdout,
	})
	c := inject.MakeCheckVars()

	in.WriteString(`---
name: go
description: go template
envs:
  PROMPT: prompt for variable
  NOPROMPT: do not prompt for variable
`)
	os.Setenv("NOPROMPT", "no prompt")

	ok, err := c.Check(in)
	if err != nil {
		t.Error(err)
	}
	assert.False(t, ok)
	assert.Contains(t, stdout.String(), `PROMPT=notset # prompt for variable`)
	assert.NotContains(t, stdout.String(), `NOPROMPT`)
}
