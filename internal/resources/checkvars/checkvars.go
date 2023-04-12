package checkvars

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kick-project/kick/internal/resources/config/configtemplate"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/logger"
	"gopkg.in/yaml.v2"
)

// Check prompt for missing variables
type Check struct {
	err    errs.HandlerIface
	log    logger.OutputIface
	stdout io.Writer
}

// Options construction options
type Options struct {
	Err    errs.HandlerIface  // Inject error hanler
	Log    logger.OutputIface // Inject logger
	Stdin  io.Reader
	Stdout io.Writer
}

// New constructor
func New(opts Options) *Check {
	p := &Check{
		err:    opts.Err,
		log:    opts.Log,
		stdout: opts.Stdout,
	}
	if p.stdout == nil {
		p.stdout = os.Stdout
	}

	return p
}

// Prompt prompt for input based on requird variables where
//
// confin is the ".kick.yml" file represented as an io.Reader
func (p *Check) Check(confin io.Reader) (bool, error) {
	buf, err := io.ReadAll(confin)
	if err != nil {
		return false, fmt.Errorf(`prompt: %w`, err)
	}

	data := &configtemplate.TemplateMain{}
	err = yaml.Unmarshal(buf, data)
	if err != nil {
		return false, fmt.Errorf(`prompt: %w`, err)
	}
	miss := map[string]string{}
	for k, v := range data.Envs {
		if os.Getenv(k) == "" {
			miss[k] = v
		}
	}
	if len(miss) > 0 {
		fmt.Fprintf(p.stdout, "## Required variables. Add these to \"%s\" file or set as environment variables.\n",
			filepath.Join(os.Getenv("HOME"), ".env"))
		for k, v := range miss {
			fmt.Fprintf(p.stdout, "%s=notset # %s\n", k, v)
		}
		return false, nil
	}

	return true, nil
}
