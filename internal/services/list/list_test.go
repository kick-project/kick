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
	assert.Regexp(t, `\|\s+HANDLE\s+\|\s+DESCRIPTION\s+\|\s+TEMPLATE\s+\|\s+LOCATION\s+\|`, out)
	assert.Regexp(t, `\|\s+handle1\s+\|\s+-\s+|\s+template1/origin1\s+\|\s+http://\S+`, out)
	assert.Regexp(t, `\|\s+handle2\s+\|\s+-\s+|\s+template2/origin1\s+\|\s+http://\S+`, out)
	assert.Regexp(t, `\|\s+handle3\s+\|\s+-\s+\|\s+template3\s+\|\s+http://\S+`, out)
	assert.Regexp(t, `\|\s+handle4\s+\|\s+-\s+\|\s+-\s+\|\s+http://\S+`, out)
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
