package templates

import (
	"fmt"
	"github.com/dengpju/higo-utils/utils"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

type Model struct {
	DB            *gorm.DB
	database      string
	prefix        string
	TableName     string
	Package       string
	Dir           string
	ModelImpl     string
	Fields        []Field
	TplFields     []TplField
}

func NewModel(DB *gorm.DB, name, dir, db, pre string) *Model {
	pkg := generator.CamelCase(strings.TrimLeft(name, pre)) + "Model"
	return &Model{DB: DB, TableName: name, Package: pkg, ModelImpl: "ModelImpl", Dir: dir + utils.PathSeparator() + pkg, database: db, prefix: pre}
}

func (this *Model) Template(tplfile string) string {
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

func (this *Model) Generate() {
	this.Fields = this.GetFields(this.TableName)
	for _, f := range this.Fields {
		tField := TplField{
			Field:   generator.CamelCase(f.Field),
			Type:    getFiledType(f),
			DbField: f.Field,
		}
		if f.Key == "PRI" {
			tField.Field = "ID"
		}
		this.TplFields = append(this.TplFields, tField)
	}
	// 目录不存在，并创建
	if _, err := os.Stat(this.Dir); os.IsNotExist(err) {
		if os.Mkdir(this.Dir, os.ModePerm) != nil {
		}
	}
	outfile := utils.File{Name: this.Dir + utils.PathSeparator() + "model.go"}
	if outfile.Exist() {
		log.Fatalln(outfile.Name + " already existed")
	}
	modelFile, err := os.OpenFile(outfile.Name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer modelFile.Close()
	tpl := this.Template("model.tpl")
	tmpl, err := template.New(model).Parse(tpl)
	if err != nil {
		panic(err)
	}
	//生成model.go
	err = tmpl.Execute(modelFile, this)
	if err != nil {
		panic(err)
	}

	outfile = utils.File{Name: this.Dir + utils.PathSeparator() + "attributes.go"}
	if outfile.Exist() {
		log.Fatalln(outfile.Name + " already existed")
	}
	attributesFile, err := os.OpenFile(outfile.Name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer attributesFile.Close()
	tpl = this.Template("attributes.tpl")
	tmpl, err = template.New(attributes).Parse(tpl)
	if err != nil {
		panic(err)
	}
	//生成attributes.go
	err = tmpl.Execute(attributesFile, this)
	if err != nil {
		panic(err)
	}
	fmt.Println("generate model success!")
}

func (this *Model) GetTables(tableNames ...string) []Table {
	db := this.DB
	var tables []Table
	if len(tableNames) == 0 {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" + this.database + "';").Find(&tables)
	} else {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE TABLE_NAME IN (" + strings.Join(tableNames, ",") + ") AND table_schema='" + this.database + "';").Find(&tables)
	}
	return tables
}

//获取所有字段信息
func (this *Model) GetFields(tableName string) []Field {
	db := this.DB
	var fields []Field
	db.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	return fields
}

type Table struct {
	Name    string `gorm:"column:Name"`
	Comment string `gorm:"column:Comment"`
}

type TplField struct {
	Field   string
	Type    string
	DbField string
}

type Field struct {
	Field      string `gorm:"column:Field"`
	Type       string `gorm:"column:Type"`
	Null       string `gorm:"column:Null"`
	Key        string `gorm:"column:Key"`
	Default    string `gorm:"column:Default"`
	Extra      string `gorm:"column:Extra"`
	Privileges string `gorm:"column:Privileges"`
	Comment    string `gorm:"column:Comment"`
}

//获取字段类型
func getFiledType(field Field) string {
	types := strings.Split(field.Type, "(")
	switch types[0] {
	case "int":
		return "int"
	case "integer":
		return "int"
	case "mediumint":
		return "int"
	case "bit":
		return "int"
	case "year":
		return "int"
	case "smallint":
		return "int"
	case "tinyint":
		return "int"
	case "bigint":
		return "int64"
	case "decimal":
		return "float32"
	case "double":
		return "float32"
	case "float":
		return "float32"
	case "real":
		return "float32"
	case "numeric":
		return "float32"
	case "timestamp":
		return "time.Time"
	case "datetime":
		return "time.Time"
	case "time":
		return "time.Time"
	default:
		return "string"
	}
}
