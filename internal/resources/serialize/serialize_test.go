package serialize_test

import (
	"testing"

	"github.com/kick-project/kick/internal/resources/config/configtemplate"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestSerialize(t *testing.T) {
	tmpl := &configtemplate.TemplateMain{}
	txt := `---
name: goms
description: Go micro services template
envs:
  GOSERVER: E.G. Github
  GOGROUP: Go group
`

	err := yaml.Unmarshal([]byte(txt), tmpl)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, `goms`, tmpl.Name)
	assert.Equal(t, `Go micro services template`, tmpl.Desc)

	data := map[string]string{
		"GOSERVER": "E.G. Github",
		"GOGROUP":  "Go group",
	}
	assert.Equal(t, data, tmpl.Envs)
}
