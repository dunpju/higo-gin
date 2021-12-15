package templates

import (
	"fmt"
	"github.com/dengpju/higo-utils/utils"
	"github.com/dengpju/higo-utils/utils/fileutils"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"text/template"
)

type Param struct {
	Package    string
	StructName string
	OutDir     string
	FileName   string
}

func NewParam(name string, out string) *Param {
	humpUnpreName := generator.CamelCase(name)
	pkg := "Param" + humpUnpreName
	outDir := out + utils.PathSeparator() + pkg
	file := pkg + ".go"
	return &Param{Package: pkg, StructName: humpUnpreName, OutDir: outDir, FileName: file}
}

func (this *Param) Template(tplfile string) string {
	_, file, _, _ := runtime.Caller(0)
	file = path.Dir(file) + utils.PathSeparator() + tplfile
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

func (this *Param) Generate() {
	if fileutils.FileExist(this.FileName) {
		log.Fatalln(this.FileName + " already existed")
	}
	if _, err := os.Stat(this.OutDir); os.IsNotExist(err) {
		if err = os.Mkdir(this.OutDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	tpl := this.Template("param.tpl")
	tmpl, err := template.New("param.tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	outfile := fileutils.File{Name: this.OutDir + utils.PathSeparator() + this.FileName}
	paramFile := fileutils.NewFile(outfile.Name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if !paramFile.Exist() {
		paramFile.Create()
	}
	defer paramFile.Close()
	//生成param.go
	err = tmpl.Execute(paramFile, this)
	if err != nil {
		panic(err)
	}
	fmt.Println("param: " + this.OutDir + " generate success!")
}
