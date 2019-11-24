package internal

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	home, _ := filepath.Abs("../tmp/home")
	conf := LoadConfig(home, ".prjstart.yml")
	assert.NotNil(t, conf)
}
