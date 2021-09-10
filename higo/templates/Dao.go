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

type DaoMap struct {
	Key   string
	Value interface{}
	Doc   string
}

func NewDaoMap(key string, value interface{}, doc string) *DaoMap {
	return &DaoMap{Key: utils.Ucfirst(utils.CaseToCamel(key)), Value: value, Doc: doc}
}

type Dao struct {
	PackageName           string
	Imports               map[string]string
	StructName            string
	ModelPackageName      string
	ModelName             string
	EntityPackageName     string
	EntityName            string
	PrimaryId             string //大驼峰
	SmallHumpPrimaryId    string //小驼峰
	PrimaryIdType         string
	TablePrimaryId        string
	TableFields           []TableField
	ModelFields           []StructField
	LenModelFields        int
	HasCreateTime         bool
	HasUpdateTime         bool
	HasDeleteTime         bool
	EntityUpdateTimeField string
	OutDir                string
	FileName              string
}

const (
	DaoStructName = "Dao"
	DaoDirSuffix  = "Dao"
	DaoFileName   = "dao"
)

func NewDao(modelTool ModelTool, model Model, entity Entity) *Dao {
	packageName := model.HumpUnpreTableName + DaoDirSuffix
	modName := GetModName() + utils.PathSeparator()
	modelImport := `"` + modName + model.OutDir + `"`
	modelImport = strings.ReplaceAll(modelImport, utils.PathSeparator(), "/")
	entityImport := `"` + modName + entity.OutDir + `"`
	entityImport = strings.ReplaceAll(entityImport, utils.PathSeparator(), "/")
	return &Dao{
		PackageName: packageName,
		Imports: map[string]string{
			"modelImport":  modelImport,
			"entityImport": entityImport,
		},
		StructName:            DaoStructName,
		ModelPackageName:      model.PackageName,
		ModelName:             model.StructName,
		EntityPackageName:     entity.PackageName,
		EntityName:            entity.StructName,
		PrimaryId:             model.PrimaryId,
		SmallHumpPrimaryId:    model.SmallHumpPrimaryId,
		PrimaryIdType:         model.PrimaryIdType,
		TablePrimaryId:        model.TablePrimaryId,
		TableFields:           model.TableFields,
		ModelFields:           model.StructFields,
		LenModelFields:        len(model.StructFields) - 1,
		HasCreateTime:         model.HasCreateTime,
		HasUpdateTime:         model.HasUpdateTime,
		HasDeleteTime:         model.HasDeleteTime,
		EntityUpdateTimeField: entity.UpdateTimeField,
		OutDir:                modelTool.OutDaoDir + utils.PathSeparator() + packageName,
		FileName:              DaoFileName,
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
	fileName := this.OutDir + utils.PathSeparator() + this.FileName + ".go"
	if utils.FileExist(fileName) {
		log.Println(fileName + " already existed")
		return
	}
	outFile := utils.NewFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	defer outFile.Close()
	tpl := this.Template(DaoFileName + ".tpl")
	tmpl, err := template.New(DaoFileName + ".tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	//生成dao.go
	err = tmpl.Execute(outFile.File(), this)
	if err != nil {
		panic(err)
	}
	fmt.Println("dao: " + this.OutDir + " generate success!")
}
