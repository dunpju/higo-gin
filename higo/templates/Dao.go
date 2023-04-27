package templates

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo/templates/tpls"
	"github.com/dengpju/higo-utils/utils/dirutil"
	"github.com/dengpju/higo-utils/utils/fileutil"
	"github.com/dengpju/higo-utils/utils/stringutil"
	"os"
	"strings"
)

type DaoMap struct {
	Key   string
	Value interface{}
	Doc   string
}

func NewDaoMap(key string, value interface{}, doc string) *DaoMap {
	return &DaoMap{Key: stringutil.Ucfirst(stringutil.CaseToCamel(key)), Value: value, Doc: doc}
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
	ModelEndField         string
	HasCreateTime         bool
	HasUpdateTime         bool
	HasDeleteTime         bool
	EntityUpdateTimeField string
	EntityDeleteTimeField string
	OutDir                string
	FileName              string
}

const (
	DaoStructName = "Dao"
	DaoDirSuffix  = "Dao"
	DaoFileName   = "dao"
)

func NewDao(modelTool ModelTool, model Model, entity Entity) *Dao {
	packageName := dirutil.Basename(modelTool.OutDaoDir)
	daoStructName := model.HumpUnpreTableName + DaoDirSuffix
	daoFileName := daoStructName
	modName := GetModName() + dirutil.PathSeparator()
	modelImport := `"` + modName + model.OutDir + `"`
	modelImport = strings.ReplaceAll(modelImport, dirutil.PathSeparator(), "/")
	entityImport := `"` + modName + entity.OutDir + `"`
	entityImport = strings.ReplaceAll(entityImport, dirutil.PathSeparator(), "/")
	return &Dao{
		PackageName: packageName,
		Imports: map[string]string{
			"modelImport":  modelImport,
			"entityImport": entityImport,
		},
		StructName:            daoStructName,
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
		ModelEndField:         model.EndField,
		HasCreateTime:         model.HasCreateTime,
		HasUpdateTime:         model.HasUpdateTime,
		HasDeleteTime:         model.HasDeleteTime,
		EntityUpdateTimeField: entity.UpdateTimeField,
		EntityDeleteTimeField: entity.DeleteTimeField,
		OutDir:                modelTool.OutDaoDir,
		FileName:              daoFileName,
	}
}

func (this *Dao) Template(tplfile string) *tpls.Tpl {
	return tpls.New(tplfile)
}

func (this *Dao) Generate() {
	this.generate()
}

func (this *Dao) generate() {
	dirutil.Dir(this.OutDir).Create()
	fileName := this.OutDir + dirutil.PathSeparator() + this.FileName + ".go"
	if fileutil.FileExist(fileName) {
		fmt.Println("dao: " + fileName + " already existed")
		return
	}
	outFile := fileutil.NewFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if !outFile.Exist() {
		outFile.Create()
	}
	defer outFile.Close()
	tmpl, err := this.Template(DaoFileName + ".tpl").Parse()
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
