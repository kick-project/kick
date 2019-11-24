package internal

import (
	"fmt"
	"github.com/crosseyed/prjstart/internal/template"
	"path/filepath"
	"strconv"
	"time"
)

type PrjVars [][]string

func SetVars(opts *OptStart) PrjVars {
	curTime := time.Now()
	data := [][]string{
		{"NAME", filepath.Base(opts.Project), "Project Name derived from target directory"},
		{"YEAR", strconv.Itoa(curTime.Year()), "Current year"},
		{"MON", strconv.Itoa(int(curTime.Month())), "Current month"},
		{"DAY", strconv.Itoa(curTime.Day()), "Current day"},
		{"HOUR", strconv.Itoa(curTime.Hour()), "Current hour"},
		{"MIN", strconv.Itoa(curTime.Minute()), "Current minute"},
		{"SEC", strconv.Itoa(curTime.Second()), "Current second"},
		{"PMON", fmt.Sprintf("%02d", int(curTime.Month())), "Zero padded month"},
		{"PDAY", fmt.Sprintf("%02d", curTime.Day()), "Zero padded day"},
		{"PHOUR", fmt.Sprintf("%02d", curTime.Hour()), "Zero padded hour"},
		{"PMIN", fmt.Sprintf("%02d", curTime.Minute()), "Zero padded minute"},
		{"PSEC", fmt.Sprintf("%02d", curTime.Second()), "Zero padded second"},
	}
	return data
}

func (prjvars PrjVars) GetVars() *template.TmplVars {
	data := map[string]string{}
	for _, v := range prjvars {
		data[v[0]] = v[1]
	}
	return template.NewTmplVars(data)
}

func (prjvars PrjVars) GetDescriptions() [][]string {
	data := [][]string{}
	for _, v := range prjvars {
		data = append(data, []string{v[0], v[2]})
	}
	return data
}
