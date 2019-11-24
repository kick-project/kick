package internal

import (
	"log"
	"path/filepath"
)

func init() {
	home, _ := filepath.Abs("../tmp/home")
	Config = LoadConfig(home, "")
	if Config == nil {
		log.Fatalf("Config is nil\n\tHOME: %s\n", home)
	}
}
