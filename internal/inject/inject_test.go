package inject

import (
	"os"
	"testing"

	"github.com/crosseyed/prjstart/internal/utils/dfaults"
	"syreclabs.com/go/faker"
)

func TestMain(m *testing.M) {
	Init()
	ret := m.Run()
	Reset()
	os.Exit(ret)
}

func TestSet(t *testing.T) {
	defaultText := "default string"
	injectText := faker.Lorem().Word()
	Set("testset_handle", injectText)
	data, ok := dfaults.Interface(defaultText, Get("testset_handle")).(string)
	if !ok || data != injectText {
		t.Fail()
	}

	Reset()
	data, ok = dfaults.Interface(defaultText, Get("testset_handle")).(string)
	if !ok || data != defaultText {
		t.Fail()
	}
}
