package tpls

import (
	"github.com/dunpju/higo-utils/utils"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"text/template"
)

type Tpl struct {
	name    string
	context string
}

func New(name string) *Tpl {
	return &Tpl{name: name, context: context(name)}
}

func context(tplName string) string {
	_, file, _, _ := runtime.Caller(0)
	file = path.Dir(file) + utils.Dir.Separator() + tplName
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	context, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(context)
}

func (this *Tpl) Context() string {
	return this.context
}

func (this *Tpl) Parse() (*template.Template, error) {
	return template.New(this.name).Parse(this.context)
}
