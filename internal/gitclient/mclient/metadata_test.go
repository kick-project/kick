package mclient

import (
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/config"
)

func TestMetadata_Checkouts(t *testing.T) {
	fn := func(lpath string) {
	}
	home, err := filepath.Abs("../../../tmp/home")
	if err != nil {
		t.Fatalf("Can not get absoulte path: %v", err)
	}
	m := Metadata{
		Config: &config.Config{
			Home: home,
		},
	}

	_ = m
	_ = fn
}
