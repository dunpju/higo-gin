package templates

import (
	"fmt"
	"github.com/dunpju/higo-gin/higo/templates/tpls"
	"github.com/dunpju/higo-utils/utils/dirutil"
	"github.com/dunpju/higo-utils/utils/fileutil"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"log"
	"os"
)

const (
	List   = "List"
	Add    = "Add"
	Edit   = "Edit"
	Delete = "Delete"
)

type Param struct {
	Package    string
	StructName string
	OutDir     string
	FileName   string
}

func NewParam(name string, out string) *Param {
	humpUnpreName := generator.CamelCase(name)
	pkg := dirutil.Basename(out)
	outDir := out
	file := "Param" + humpUnpreName + ".go"
	return &Param{Package: pkg, StructName: humpUnpreName, OutDir: outDir, FileName: file}
}

func (this *Param) Template(tplfile string) *tpls.Tpl {
	return tpls.New(tplfile)
}

func (this *Param) Generate() {
	outFile := this.OutDir + dirutil.PathSeparator() + this.FileName
	if fileutil.FileExist(outFile) {
		log.Println(outFile + " already existed")
		return
	}
	if _, err := os.Stat(this.OutDir); os.IsNotExist(err) {
		if err = os.MkdirAll(this.OutDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	tmpl, err := this.Template("param.tpl").Parse()
	if err != nil {
		panic(err)
	}
	outfile := fileutil.File{Name: this.OutDir + dirutil.PathSeparator() + this.FileName}
	paramFile := fileutil.NewFile(outfile.Name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if !paramFile.Exist() {
		paramFile.Create()
	}
	defer paramFile.Close()
	//生成param.go
	err = tmpl.Execute(paramFile, this)
	if err != nil {
		panic(err)
	}
	fmt.Println("param: " + outfile.Name + " generate success!")
}
