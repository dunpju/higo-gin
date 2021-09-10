package templates

import (
	"fmt"
	"github.com/dengpju/higo-utils/utils"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"text/template"
)

type Entity struct {
	PackageName   string
	StructName    string
	PrimaryId     string
	StructFields  []StructField
	HasCreateTime bool
	HasUpdateTime bool
	OutDir        string
	FileName      string
}

const (
	EntityStructName   = "Impl"
	EntityDirSuffix    = "Entity"
	EntityFileName     = "entity"
	EntityFlagFileName = "flag"
)

func NewEntity(modelTool ModelTool, model Model) *Entity {
	packageName := model.HumpUnpreTableName + EntityDirSuffix
	return &Entity{
		PackageName:   packageName,
		StructName:    EntityStructName,
		PrimaryId:     model.PrimaryId,
		StructFields:  model.StructFields,
		HasCreateTime: model.HasCreateTime,
		HasUpdateTime: model.HasUpdateTime,
		OutDir:        modelTool.OutEntityDir + utils.PathSeparator() + packageName,
		FileName:      EntityFileName,
	}
}

func (this *Entity) Template(tplfile string) string {
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

func (this *Entity) Generate() {
	if _, err := os.Stat(this.OutDir); os.IsNotExist(err) {
		if err = os.Mkdir(this.OutDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	tpl := this.Template(EntityFileName + ".tpl")
	tmpl, err := template.New(EntityFileName + ".tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	outfile := utils.File{Name: this.OutDir + utils.PathSeparator() + EntityFileName + ".go"}
	modelFile, err := os.OpenFile(outfile.Name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer modelFile.Close()
	//生成entity.go
	err = tmpl.Execute(modelFile, this)
	if err != nil {
		panic(err)
	}
	outfile = utils.File{Name: this.OutDir + utils.PathSeparator() + EntityFlagFileName + ".go"}
	attributesFile, err := os.OpenFile(outfile.Name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer attributesFile.Close()
	tpl = this.Template(EntityFileName + "_" + EntityFlagFileName + ".tpl")
	tmpl, err = template.New(attributes).Parse(tpl)
	if err != nil {
		panic(err)
	}
	//生成flag.go
	err = tmpl.Execute(attributesFile, this)
	if err != nil {
		panic(err)
	}
	fmt.Println("model: " + this.OutDir + " generate success!")
}
