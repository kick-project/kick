package gitclient

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/stretchr/testify/assert"
)

//
// Tests
//
func TestGitClient_Tag(t *testing.T) {
	url := "http://127.0.0.1:5000/tmpl1.git"
	tmpdir := utils.TempDir()
	tmpproject, err := ioutil.TempDir(tmpdir, "TestGitClient_Tag-*")
	if err != nil {
		t.Fatal("Can not create temporary directory")
	}
	err = os.Remove(tmpproject)
	if err != nil {
		t.Fatal("Can not remove temporary directory")
	}
	client := Gitclient{
		URL:    url,
		Local:  tmpproject,
		Output: os.Stdout,
	}
	client.Sync()
	tags := client.Tags()
	tlen := len(tags)
	for _, tag := range tags {
		assert.Regexp(t, regexp.MustCompile(`^\d+\.\d+\.\d+$`), tag)
	}
	assert.Greater(t, tlen, 0)
}
