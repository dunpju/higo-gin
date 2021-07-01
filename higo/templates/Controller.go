package templates

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

type TplEngine interface {
	Template() string
}

type Controller struct {
	Package string
	Name    string
}

func NewController() *Controller {
	return &Controller{}
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
