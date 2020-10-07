package template

import (
	"io/ioutil"
	"os"
	tt "text/template"

	"github.com/crosseyed/prjstart/internal/services/template/variables"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

// File2File takes a src file populates a dst file with the results
// of the template populated with variables
func File2File(src, dst string, vars *variables.Variables) {
	b, err := ioutil.ReadFile(src)
	errutils.Elogf("Can not open template file %s for reading: %v", src, err)
	Txt2File(string(b), dst, vars)
}

// Txt2File takes template text tmpltxt and outputs to dst file
func Txt2File(tmpltxt, dst string, vars *variables.Variables) {
	td := os.Getenv("TEMP")
	f, err := ioutil.TempFile(td, "prjstart-*")
	errutils.Epanicf("Error creating tempfile %v", err)

	t := tt.Must(tt.New("tmpltxt").Parse(tmpltxt))

	errutils.Epanicf("Error parsing variables: %v", err)
	err = t.Execute(f, vars)
	errutils.Epanicf("Error executing template: %v", err)
	err = f.Close()
	errutils.Epanicf("Error closing tempfile: %v", err)
	err = os.Rename(f.Name(), dst)
	errutils.Epanicf("Error writing file %s: %v", dst, err)
}

// Txt2String takes a tmpltxt string and generates a string output
func Txt2String(tmpltxt string, vars *variables.Variables) string {
	td := os.Getenv("TEMP")
	f, err := ioutil.TempFile(td, "prjstart-*")
	errutils.Epanicf("Error creating tempfile %v", err)
	f.Close() // nolint
	Txt2File(tmpltxt, f.Name(), vars)
	b, err := ioutil.ReadFile(f.Name())
	errutils.Elogf("Can not open template file %s for reading: %v", f.Name(), err)
	os.Remove(f.Name()) // nolint
	return string(b)
}
