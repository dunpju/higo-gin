package templates

import (
	"fmt"
	"github.com/dengpju/higo-utils/utils"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

type Enum struct {
	Package   string
	Name      string
	OutStruct string
	File      string
	Enums     []*utils.KeyValue
}

func NewEnum(pkg string, name string, file string) *Enum {
	name = utils.Ucfirst(name)
	name = name + enum
	outStruct := file + utils.PathSeparator() + strings.Trim(name, enum) + enum
	file = outStruct + ".go"
	return &Enum{Package: pkg, Name: name, OutStruct: outStruct, File: file}
}

func (this *Enum) Template(tplfile string) string {
	_, file, _, _ := runtime.Caller(0)
	file = path.Dir(file) + utils.PathSeparator() + tplfile
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

func (this *Enum) Generate() {
	outfile := utils.File{Name: this.File}
	if outfile.Exist() {
		log.Fatalln(this.File + " already existed")
	}
	outFile, err := os.OpenFile(this.File, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	tpl := this.Template("enum.tpl")
	tmpl, err := template.New("enum.tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	//生成enum
	err = tmpl.Execute(outFile, this)
	if err != nil {
		panic(err)
	}
	fmt.Println("enum: " + this.OutStruct + " generate success!")
}
