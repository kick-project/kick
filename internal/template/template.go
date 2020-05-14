package template

import (
	"io/ioutil"
	"os"
	"text/template"

	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

// File2File takes a src file populates a dst file with the results
// of the template populated with variables
func File2File(src, dst string, vars *TmplVars) {
	b, err := ioutil.ReadFile(src)
	errutils.Elogf(err, "Can not open template file %s for reading: %v", src, err)
	Txt2File(string(b), dst, vars)
}

// Txt2File takes template text tmpltxt and outputs to dst file
func Txt2File(tmpltxt, dst string, vars *TmplVars) {
	td := os.Getenv("TEMP")
	f, err := ioutil.TempFile(td, "prjstart-*")
	errutils.Epanicf(err, "Error creating tempfile %v", err)

	t := template.Must(template.New("tmpltxt").Parse(tmpltxt))

	errutils.Epanicf(err, "Error parsing variables: %v", err)
	err = t.Execute(f, vars)
	errutils.Epanicf(err, "Error executing template: %v", err)
	err = f.Close()
	errutils.Epanicf(err, "Error closing tempfile: %v", err)
	err = os.Rename(f.Name(), dst)
	errutils.Epanicf(err, "Error writing file %s: %v", dst, err)
}

// Txt2String takes a tmpltxt string and generates a string output
func Txt2String(tmpltxt string, vars *TmplVars) string {
	td := os.Getenv("TEMP")
	f, err := ioutil.TempFile(td, "prjstart-*")
	errutils.Epanicf(err, "Error creating tempfile %v", err)
	f.Close() // nolint
	Txt2File(tmpltxt, f.Name(), vars)
	b, err := ioutil.ReadFile(f.Name())
	errutils.Elogf(err, "Can not open template file %s for reading: %v", f.Name(), err)
	os.Remove(f.Name()) // nolint
	return string(b)
}
