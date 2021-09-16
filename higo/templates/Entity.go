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
	PackageName     string
	Imports         map[string]string
	StructName      string
	PrimaryId       string
	StructFields    []StructField
	HasCreateTime   bool
	HasUpdateTime   bool
	UpdateTimeField string
	DeleteTimeField string
	OutDir          string
	FileName        string
}

const (
	EntityStructName   = "Impl"
	EntityDirSuffix    = "Entity"
	EntityFileName     = "entity"
	EntityFlagFileName = "flag"
)

func NewEntity(modelTool ModelTool, model Model) *Entity {
	packageName := model.HumpUnpreTableName + EntityDirSuffix
	imports := make(map[string]string)
	if model.HasCreateTime || model.HasUpdateTime {
		imports["time"] = `"time"`
	}
	return &Entity{
		PackageName:     packageName,
		Imports:         imports,
		StructName:      EntityStructName,
		PrimaryId:       model.PrimaryId,
		StructFields:    model.StructFields,
		HasCreateTime:   model.HasCreateTime,
		HasUpdateTime:   model.HasUpdateTime,
		UpdateTimeField: model.UpdateTimeField,
		DeleteTimeField: model.DeleteTimeField,
		OutDir:          modelTool.OutEntityDir + utils.PathSeparator() + packageName,
		FileName:        EntityFileName,
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
	entityFile, err := os.OpenFile(outfile.Name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer entityFile.Close()
	//生成entity.go
	err = tmpl.Execute(entityFile, this)
	if err != nil {
		panic(err)
	}
	outFlagFile := utils.File{Name: this.OutDir + utils.PathSeparator() + EntityFlagFileName + ".go"}
	if !utils.FileExist(outFlagFile.Name) { // flag.go 不存在则生成
		flagFile, err := os.OpenFile(outFlagFile.Name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer flagFile.Close()
		tpl = this.Template(EntityFileName + "_" + EntityFlagFileName + ".tpl")
		tmpl, err = template.New(attributes).Parse(tpl)
		if err != nil {
			panic(err)
		}
		//生成flag.go
		err = tmpl.Execute(flagFile, this)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("entity: " + this.OutDir + " generate success!")
}