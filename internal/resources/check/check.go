// Package check performs checks
package check

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/logger"
)

//go:generate ifacemaker -f check.go -s Check -p check -i CheckIface -o check_interfaces.go -c "AUTO GENERATED"

// Check runs a series of checks
type Check struct {
	configPath         string             `validate:"required"`
	configTemplatePath string             `validate:"required"`
	home               string             `validate:"required"`
	err                *errs.Handler      `validate:"required"`
	log                logger.OutputIface `validate:"required"`
	metadataDir        string             `validate:"required"`
	sqliteFile         string             `validate:"required"`
	stderr             io.Writer          `validate:"required"`
	stdout             io.Writer          `validate:"required"`
	templateDir        string             `validate:"required"`
}

// Options options for constructor
type Options struct {
	ConfigPath         string             `validate:"required"`
	ConfigTemplatePath string             `validate:"required"`
	Err                *errs.Handler      `validate:"required"`
	HomeDir            string             `validate:"required"`
	Log                logger.OutputIface `validate:"required"`
	MetadataDir        string             `validate:"required"`
	SQLiteFile         string             `validate:"required"`
	Stderr             io.Writer          `validate:"required"`
	Stdout             io.Writer          `validate:"required"`
	TemplateDir        string             `validate:"required"`
}

// New constructor
func New(opts *Options) *Check {
	return &Check{
		configPath:         opts.ConfigPath,
		configTemplatePath: opts.ConfigTemplatePath,
		err:                opts.Err,
		home:               opts.HomeDir,
		log:                opts.Log,
		metadataDir:        opts.MetadataDir,
		sqliteFile:         opts.SQLiteFile,
		stderr:             opts.Stderr,
		stdout:             opts.Stdout,
		templateDir:        opts.TemplateDir,
	}
}

// Init checks to see if an initialization has been performed. This function
// will print an error message and exit if initialization is needed.
func (c *Check) Init() error {
	msg := "not initialized. please run \"kick setup\" to initialize configuration"
	// Directory checks
	dirs := []string{c.home, c.metadataDir, c.templateDir}
	for _, d := range dirs {
		info, err := os.Stat(d)
		if os.IsNotExist(err) {
			return errors.New(msg)
		}
		c.err.Panic(err)

		if !info.IsDir() {
			return fmt.Errorf("warning %s is not a directory. please remove then run \"kick init\" to initialize", d)
		}
	}

	// File checks
	files := []string{c.configPath, c.configTemplatePath, c.sqliteFile}
	for _, f := range files {
		info, err := os.Stat(f)
		if os.IsNotExist(err) {
			return errors.New(msg)
		}
		c.err.Panic(err)

		if info.IsDir() {
			return fmt.Errorf("expected a normal file %s got a directory. please remove then run \"kick init\" to initialize", f)
		}
	}

	return nil
}
