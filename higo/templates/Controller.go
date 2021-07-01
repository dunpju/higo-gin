package templates

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"text/template"
)

type TplEngine interface {
	Template() string
	Generate()
}

type Controller struct {
	Package string
	Name    string
	File    string
}

func NewController(pak string, name string, file string) *Controller {
	return &Controller{Package: pak, Name: name, File: file}
}

func (this *Controller) Template() string {
	_, file, _, _ := runtime.Caller(0)
	file = strings.TrimRight(file, ".go") + ".tpl"
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	context, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(context)
}

func (this *Controller) Generate() {
	fi, err := os.OpenFile(this.File, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	tpl := this.Template()
	tmpl, err := template.New("Controller").Parse(tpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(fi, this)
	if err != nil {
		panic(err)
	}
}
