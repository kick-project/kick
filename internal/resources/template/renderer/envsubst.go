package renderer

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/drone/envsubst"
	"github.com/kick-project/kick/internal/resources/file"
	"github.com/kick-project/kick/internal/resources/template/variables"
)

// RenderEnv renders using environment variables
type RenderEnv struct {
	Renderer
}

// File2File takes a src file populates a dst file with the results of the
// template populated with variables
func (r *RenderEnv) File2File(src, dst string, vars *variables.Variables, nounset, noempty bool) (err error) {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return
	}
	err = r.Text2File(string(b), dst, vars, nounset, noempty)
	return
}

// Text2File takes template text text and outputs to dst file
func (r *RenderEnv) Text2File(text, dst string, vars *variables.Variables, nounset, noempty bool) (err error) {
	td := os.Getenv("TEMP")
	f, err := ioutil.TempFile(td, "kick-*")
	if err != nil {
		return fmt.Errorf("Text2File: %w", err)
	}

	result, err := r.Text2String(text, vars, nounset, noempty)
	if err != nil {
		return fmt.Errorf("Text2File: %w", err)
	}

	_, err = f.WriteString(result)
	if err != nil {
		return fmt.Errorf("Text2File: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("Text2File: %w", err)
	}
	err = file.MoveAll(f.Name(), dst)
	if err != nil {
		return fmt.Errorf("Text2File: %w", err)
	}
	return
}

// Text2String renders input text and returns result as a string.
func (r *RenderEnv) Text2String(text string, vars *variables.Variables, nounset, noempty bool) (result string, err error) {
	result, err = envsubst.EvalEnv(text)
	return
}

// RenderDirRegexp returns the regex to match directory names that should be rendered.
func (r *RenderEnv) RenderDirRegexp() *regexp.Regexp {
	regex := regexp.MustCompile(`\${[A-Za-z0-9_]+}`)
	return regex
}
