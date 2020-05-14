package utils

import (
	"os"
	"path/filepath"
	"time"
)

// BaseProjectPath returns the path to where all git projects are synced
func BaseProjectPath(home string) string {
	if home == "" {
		home = os.Getenv("HOME")
	}
	p := filepath.Join(home, ".prjstart", "projects")
	return p
}

// BaseMetadataPath returns the path to where all metadata is synced
func BaseMetadataPath(home string) string {
	if home == "" {
		home = os.Getenv("HOME")
	}
	p := filepath.Join(home, ".prjstart", "metadata")
	return p
}

// Touch touches a file
func Touch(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(filename, currentTime, currentTime)
		if err != nil {
			return err
		}
	}
	return nil
}
