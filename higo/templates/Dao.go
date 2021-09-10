package templates

import (
	"fmt"
	"github.com/dengpju/higo-utils/utils"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

type DaoMap struct {
	Key   string
	Value interface{}
	Doc   string
}

func NewDaoMap(key string, value interface{}, doc string) *DaoMap {
	return &DaoMap{Key: utils.Ucfirst(utils.CaseToCamel(key)), Value: value, Doc: doc}
}

type Dao struct {
	PackageName        string
	Imports            map[string]string
	StructName         string
	ModelPackageName   string
	ModelName          string
	EntityPackageName  string
	EntityName         string
	PrimaryId          string //大驼峰
	SmallHumpPrimaryId string //小驼峰
	PrimaryIdType      string
	TablePrimaryId     string
	TableFields        []TableField
	ModelFields        []StructField
	HasDeleteTime      bool
	OutDir             string
	FileName           string
}

const (
	DaoStructName = "Dao"
	DaoDirSuffix  = "Dao"
	DaoFileName   = "dao"
)

func NewDao(modelTool ModelTool, model Model) *Dao {
	packageName := generator.CamelCase(strings.Replace(model.TableName, model.Prefix, "", 1)) + DaoDirSuffix
	return &Dao{
		PackageName: packageName,
		Imports:            make(map[string]string),
		StructName:         DaoStructName,
		ModelPackageName:   model.PackageName,
		ModelName:          model.StructName,
		EntityPackageName:  model.StructName,
		EntityName:         model.StructName,
		PrimaryId:          model.PrimaryId,
		SmallHumpPrimaryId: model.SmallHumpPrimaryId,
		PrimaryIdType:      model.PrimaryIdType,
		TablePrimaryId:     model.TablePrimaryId,
		TableFields:        model.TableFields,
		ModelFields:        model.StructFields,
		HasDeleteTime:      model.HasDeleteTime,
		OutDir:             modelTool.OutDaoDir + utils.PathSeparator() + packageName,
		FileName:           DaoFileName,
	}
}

func (this *Dao) Template(tplfile string) string {
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

func (this *Dao) Generate() {
	this.generate()
}

func (this *Dao) generate() {
	utils.Dir(this.OutDir).Create()
	if utils.FileExist(this.FileName) {
		log.Println(this.FileName + " already existed")
		return
	}
	outFile := utils.NewFile(this.FileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	defer outFile.Close()
	tpl := this.Template("dao.tpl")
	tmpl, err := template.New("dao.tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	//生成enum
	err = tmpl.Execute(outFile.File(), this)
	if err != nil {
		panic(err)
	}
	fmt.Println("dao: " + this.StructName + " generate success!")
}
