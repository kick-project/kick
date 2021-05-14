// Package check performs checks
package check

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kick-project/kick/internal/resources/errs"
)

// Check runs a series of checks
type Check struct {
	ConfigPath         string    `validate:"required"`
	ConfigTemplatePath string    `validate:"required"`
	HomeDir            string    `validate:"required"`
	MetadataDir        string    `validate:"required"`
	SQLiteFile         string    `validate:"required"`
	Stderr             io.Writer `validate:"required"`
	Stdout             io.Writer `validate:"required"`
	TemplateDir        string    `validate:"required"`
}

// Init checks to see if an initialization has been performed. This function
// will print an error message and exit if initialization is needed.
func (c *Check) Init() error {
	msg := "not initialized. please run \"kick setup\" to initialize configuration"
	// Directory checks
	dirs := []string{c.HomeDir, c.MetadataDir, c.TemplateDir}
	for _, d := range dirs {
		info, err := os.Stat(d)
		if os.IsNotExist(err) {
			return errors.New(msg)
		}
		errs.Panic(err)

		if !info.IsDir() {
			return fmt.Errorf("warning %s is not a directory. please remove then run \"kick init\" to initialize", d)
		}
	}

	// File checks
	files := []string{c.ConfigPath, c.ConfigTemplatePath, c.SQLiteFile}
	for _, f := range files {
		info, err := os.Stat(f)
		if os.IsNotExist(err) {
			return errors.New(msg)
		}
		errs.Panic(err)

		if info.IsDir() {
			return fmt.Errorf("expected a normal file %s got a directory. please remove then run \"kick init\" to initialize", f)
		}
	}

	return nil
}
