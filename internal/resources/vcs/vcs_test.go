package vcs_test

import (
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/kick-project/kick/internal/resources/vcs"
	"github.com/stretchr/testify/assert"
)

var params *testParams

type testParams struct {
	base     string
	di       *di.DI
	repopath string
	vcs      *vcs.VCS
}

func setup() *testParams {
	if params != nil {
		return params
	}
	params = &testParams{}
	params.base = filepath.Join(testtools.TempDir(), "TestInfo")
	params.repopath = filepath.Join(params.base, "testrepo")
	params.di = di.Setup(params.base)
	params.vcs = params.di.MakeVCS()
	return params
}

func TestRepo_Checkout(t *testing.T) {
	params = setup()
	r, err := params.vcs.Open(params.repopath)
	if err != nil {
		t.Error(err)
	}
	err = r.Checkout(`1.0.0`)
	if err != nil {
		t.Error(err)
	}
	err = r.Checkout(`master`)
	if err != nil {
		t.Error(err)
	}
}

func TestRepo_Versions(t *testing.T) {
	params := setup()
	r, err := params.vcs.Open(params.repopath)
	if err != nil {
		t.Error(err)
	}
	versions := r.Versions()
	fnVerCheck := func(target string) bool {
		for _, v := range versions {
			if target == v {
				return true
			}
		}
		return false
	}
	checkVersions := []string{"1.0.0", "1.1.0", "2.0.0", "2.1.0", "2.1.1"}
	for _, v := range checkVersions {
		if !fnVerCheck(v) {
			t.Errorf(`could not find version: %s`, v)
		}
	}
	assert.Equal(t, len(versions), len(checkVersions))
}
