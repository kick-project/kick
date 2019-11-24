package template

import (
	"github.com/crosseyed/prjstart/internal/utils"
	"io/ioutil"
	"os"
	"text/template"
)

// TmplFile2File takes a src file populates a dst file with the results
// of the template populated with variables
func TmplFile2File(src, dst string, vars *TmplVars) {
	b, err := ioutil.ReadFile(src)
	utils.ChkErr(err, utils.Elogf, "Can not open template file %s for reading: %v", src, err)
	TmplTxt2File(string(b), dst, vars)
}

// TmplTxt2File takes template text tmpltxt and outputs to dst file
func TmplTxt2File(tmpltxt, dst string, vars *TmplVars) {
	td := os.Getenv("TEMP")
	f, err := ioutil.TempFile(td, "prjstart-*")
	utils.ChkErr(err, utils.Epanicf, "Error creating tempfile %v", err)

	t := template.Must(template.New("tmpltxt").Parse(tmpltxt))

	utils.ChkErr(err, utils.Epanicf, "Error parsing variables: %v", err)
	err = t.Execute(f, vars)
	utils.ChkErr(err, utils.Epanicf, "Error executing template: %v", err)
	err = f.Close()
	utils.ChkErr(err, utils.Epanicf, "Error closing tempfile: %v", err)
	err = os.Rename(f.Name(), dst)
	utils.ChkErr(err, utils.Epanicf, "Error writing file %s: %v", dst, err)
}

// TmplTxt2String takes a tmpltxt string and generates a string output
func TmplTxt2String(tmpltxt string, vars *TmplVars) string {
	td := os.Getenv("TEMP")
	f, err := ioutil.TempFile(td, "prjstart-*")
	utils.ChkErr(err, utils.Epanicf, "Error creating tempfile %v", err)
	f.Close() // nolint
	TmplTxt2File(tmpltxt, f.Name(), vars)
	b, err := ioutil.ReadFile(f.Name())
	utils.ChkErr(err, utils.Elogf, "Can not open template file %s for reading: %v", f.Name(), err)
	os.Remove(f.Name()) // nolint
	return string(b)
}
