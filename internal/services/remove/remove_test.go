package remove_test

import (
	"path/filepath"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/iinitialize"
	"github.com/kick-project/kick/internal/di/iremove"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/services/remove"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/file"
	"github.com/stretchr/testify/assert"
)

func TestRemoveFirst(t *testing.T) {
	r := setup()
	r.Remove("handle1")
	assert.NotContains(t, tList(r.Conf), "handle1")
	assert.Contains(t, tList(r.Conf), "handle2")
	assert.Contains(t, tList(r.Conf), "handle3")
}

func TestRemoveMiddle(t *testing.T) {
	r := setup()
	r.Remove("handle2")
	assert.Contains(t, tList(r.Conf), "handle1")
	assert.NotContains(t, tList(r.Conf), "handle2")
	assert.Contains(t, tList(r.Conf), "handle3")
}

func TestRemoveLast(t *testing.T) {
	r := setup()
	r.Remove("handle3")
	assert.Contains(t, tList(r.Conf), "handle1")
	assert.Contains(t, tList(r.Conf), "handle2")
	assert.NotContains(t, tList(r.Conf), "handle3")
}

func setup() (r *remove.Remove) {
	home := filepath.Clean(utils.TempDir() + "/TestRemove")

	src := filepath.Clean(home + "/.kick/templates.yml.save")
	dst := filepath.Clean(home + "/.kick/templates.yml")
	_, err := file.Copy(src, dst)
	if err != nil {
		panic(err)
	}
	inject := di.Setup(home)

	init := &initialize.Initialize{}
	err = copier.Copy(init, iinitialize.Inject(inject))
	if err != nil {
		panic(err)
	}

	r = &remove.Remove{}
	err = copier.Copy(r, iremove.Inject(inject))
	errutils.Epanic(err)
	return r
}

func tList(temp *config.File) (list []string) {
	for _, t := range temp.Templates {
		list = append(list, t.Handle)
	}
	return
}
