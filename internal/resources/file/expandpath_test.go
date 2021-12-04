package file_test

import (
	"testing"

	"github.com/kick-project/kick/internal/resources/file"
	"github.com/stretchr/testify/assert"
)

// TestExpandPath test path expansion
func TestExpandPath(t *testing.T) {
	original := "~/.bashrc"
	expanded := file.ExpandPath(original)
	assert.NotContains(t, expanded, "~")
	assert.Greater(t, len(expanded), len(original))
}
