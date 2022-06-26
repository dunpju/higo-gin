package templates

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo/templates/tpls"
	"github.com/dengpju/higo-utils/utils/dirutil"
	"github.com/dengpju/higo-utils/utils/fileutil"
	"github.com/dengpju/higo-utils/utils/stringutil"
	"log"
	"os"
	"strings"
)

type Class struct {
	Package   string
	Name      string
	OutDir    string
	OutStruct string
	SelfName  string
	File      string
}

func NewClass(pkg string, name string, file string) *Class {
	name = stringutil.Ucfirst(name)
	outStruct := file + dirutil.PathSeparator() + name
	c := &Class{Package: pkg, Name: name, OutStruct: outStruct, SelfName: strings.ReplaceAll(outStruct, "\\", "/")}
	c.OutDir = file
	file = outStruct + ".go"
	c.File = file
	return c
}

func (this *Class) Template(tplfile string) *tpls.Tpl {
	return tpls.New(tplfile)
}

func (this *Class) Generate() {
	dirutil.Dir(this.OutDir).Create()
	if fileutil.FileExist(this.File) {
		log.Println(this.File + " already existed")
		return
	}
	outFile := fileutil.NewFile(this.File, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if !outFile.Exist() {
		outFile.Create()
	}
	defer outFile.Close()
	tmpl, err := this.Template("class.tpl").Parse()
	if err != nil {
		panic(err)
	}
	//生成class
	err = tmpl.Execute(outFile.File(), this)
	if err != nil {
		panic(err)
	}
	fmt.Println("class: " + this.OutStruct + " generate success!")
}
