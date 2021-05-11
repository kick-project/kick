package renderer

import (
	"io/ioutil"
	"os"
	"regexp"
	tt "text/template"

	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/file"
	"github.com/kick-project/kick/internal/resources/template/variables"
)

//
// RenderText
//

// RenderText render using text/template
type RenderText struct {
	Renderer
}

// File2File takes a src file populates a dst file with the results of the
// template populated with variables
func (r *RenderText) File2File(src, dst string, vars *variables.Variables, nounset, noempty bool) error {
	b, err := ioutil.ReadFile(src)
	errs.LogF("Can not open template file %s for reading: %v", src, err)
	err = r.Text2File(string(b), dst, vars, nounset, noempty)
	errs.Panic(err)
	return err
}

// Text2File takes template text text and outputs to dst file
func (r *RenderText) Text2File(text, dst string, vars *variables.Variables, nounset, noempty bool) error {
	td := os.Getenv("TEMP")
	f, err := ioutil.TempFile(td, "kick-*")
	errs.PanicF("Error creating tempfile %v", err)

	t := tt.Must(tt.New("texttemplate").Parse(text))

	errs.PanicF("Error parsing variables: %v", err)
	err = t.Execute(f, vars)
	errs.PanicF("Error executing template: %v", err)
	err = f.Close()
	errs.PanicF("Error closing tempfile: %v", err)
	err = file.Move(f.Name(), dst)
	errs.PanicF("Error writing file %s: %v", dst, err)
	return err
}

// Text2String renders input text and returns result as a string.
func (r *RenderText) Text2String(text string, vars *variables.Variables, nounset, noempty bool) (string, error) {
	td := os.Getenv("TEMP")
	f, err := ioutil.TempFile(td, "kick-*")
	errs.PanicF("Error creating tempfile %v", err)
	f.Close() // nolint
	err = r.Text2File(text, f.Name(), vars, nounset, noempty)
	errs.Panic(err)
	b, err := ioutil.ReadFile(f.Name())
	errs.LogF("Can not open template file %s for reading: %v", f.Name(), err)
	os.Remove(f.Name()) // nolint
	return string(b), err
}

// RenderDirRegexp returns the regex to match directory names that should be rendered.
func (r *RenderText) RenderDirRegexp() *regexp.Regexp {
	regex := regexp.MustCompile(`{{[^}}]+}}`)
	return regex
}
