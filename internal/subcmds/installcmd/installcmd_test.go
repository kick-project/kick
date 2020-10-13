package installcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/apex/log"
	"github.com/crosseyed/prjstart/internal/resources/db"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/crosseyed/prjstart/internal/utils/file"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", usageDoc)
}

func TestInstall(t *testing.T) {
	home := filepath.Join(utils.TempDir(), "home")

	type install struct {
		handle   string
		template string
	}

	installs := &[]install{
		{
			handle:   "handle1",
			template: "tmpl",
		},
		{
			handle:   "handle2",
			template: "tmpl1/master1",
		},
		{
			handle:   "handle3",
			template: "http://localhost:5000/tmpl2.git",
		},
	}

	s := settings.GetSettings(home)
	s.LogLevel(log.DebugLevel)
	bak := fmt.Sprintf("%s.bak", s.PathTemplateConf)
	file.Copy(s.PathTemplateConf, bak)
	defer func() {
		os.Rename(bak, s.PathTemplateConf)
		dbconn := s.GetDB()
		db.Lock()
		for _, inst := range *installs {
			_, err := dbconn.Exec(`DELETE FROM installed WHERE handle=?`, inst.handle)
			if err != nil {
				t.Error(err)
			}
		}
		db.Unlock()
	}()

	for _, inst := range *installs {
		ec := Install([]string{"install", inst.handle, inst.template}, s)
		assert.Equal(t, 0, ec)
	}
	assert.Equal(t, "", "")
}
