package templates

import (
	"fmt"
	"github.com/dunpju/higo-gin/higo/templates/tpls"
	"github.com/dunpju/higo-utils/utils/dirutil"
	"github.com/dunpju/higo-utils/utils/fileutil"
	"os"
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
		OutDir:          modelTool.OutEntityDir + dirutil.PathSeparator() + packageName,
		FileName:        EntityFileName,
	}
}

func (this *Entity) Template(tplfile string) *tpls.Tpl {
	return tpls.New(tplfile)
}

func (this *Entity) Generate() {
	if _, err := os.Stat(this.OutDir); os.IsNotExist(err) {
		if err = os.Mkdir(this.OutDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	tmpl, err := this.Template(EntityFileName + ".tpl").Parse()
	if err != nil {
		panic(err)
	}
	outfile := fileutil.File{Name: this.OutDir + dirutil.PathSeparator() + EntityFileName + ".go"}
	entityFile := fileutil.NewFile(outfile.Name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if !entityFile.Exist() {
		entityFile.Create()
	}
	defer entityFile.Close()
	//生成entity.go
	err = tmpl.Execute(entityFile, this)
	if err != nil {
		panic(err)
	}
	outFlagFile := fileutil.File{Name: this.OutDir + dirutil.PathSeparator() + EntityFlagFileName + ".go"}
	if !fileutil.FileExist(outFlagFile.Name) { // flag.go 不存在则生成
		flagFile, err := os.OpenFile(outFlagFile.Name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer flagFile.Close()
		tmpl, err = this.Template(EntityFileName + "_" + EntityFlagFileName + ".tpl").Parse()
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
