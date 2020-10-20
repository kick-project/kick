package list

import (
	"bytes"
	"testing"

	"github.com/kick-project/kick/internal/resources/config"
	"github.com/stretchr/testify/assert"
)

func TestListShort(t *testing.T) {
	stderr, stdout, conf := getOptions()
	l := List{
		Stderr: stderr,
		Stdout: stdout,
		Conf:   conf,
	}
	l.List(false)

	out := stdout.String()
	assert.Contains(t, out, "handle1")
	assert.Contains(t, out, "handle2")
}

func TestListLong(t *testing.T) {
	stderr, stdout, conf := getOptions()
	l := List{
		Stderr: stderr,
		Stdout: stdout,
		Conf:   conf,
	}
	l.List(true)
	out := stdout.String()
	assert.Regexp(t, `^Handle\s+Template\s+Location\s`, out)
	assert.Regexp(t, `handle1\s+template1/origin1\s+http://\S+`, out)
	assert.Regexp(t, `handle2\s+template2/origin1\s+http://\S+`, out)
	assert.Regexp(t, `handle3\s+template3\s+http://\S+`, out)
	assert.Regexp(t, `handle4\s+http://\S+`, out)
}

func getOptions() (stderr, stdout *bytes.Buffer, conf *config.File) {
	stderr = &bytes.Buffer{}
	stdout = &bytes.Buffer{}
	templates := []config.Template{
		{
			Handle:   "handle1",
			Template: "template1",
			Origin:   "origin1",
			URL:      "http://template.io/template1.git",
		},
		{
			Handle:   "handle2",
			Template: "template2",
			Origin:   "origin1",
			URL:      "http://template.io/template2.git",
		},
		{
			Handle:   "handle3",
			Template: "template3",
			URL:      "http://template.io/template3.git",
		},
		{
			Handle: "handle4",
			URL:    "http://template.io/template4.git",
		},
	}
	conf = &config.File{
		Stderr:    stderr,
		Templates: templates,
	}
	return
}
