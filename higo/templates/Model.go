package templates

import (
	"bytes"
	"fmt"
	"github.com/dengpju/higo-utils/utils"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/jinzhu/gorm"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

type Model struct {
	PackageName    string
	Imports        map[string]string
	StructName     string
	TableName      string
	PrimaryId      string
	PrimaryIdType  string
	TablePrimaryId string
	TableFields    []TableField
	StructFields   []StructField
	DB             *gorm.DB
	Database       string
	Prefix         string
	OutDir         string
	HasDeleteTime  bool
}

const ModelStructName = "Impl"

func NewModel(DB *gorm.DB, name, outDir, db, pre string) *Model {
	pkg := generator.CamelCase(strings.Replace(name, pre, "", 1)) + "Model"
	return &Model{DB: DB, TableName: name, PackageName: pkg, StructName: ModelStructName,
		OutDir: outDir + utils.PathSeparator() + pkg, Database: db, Prefix: pre,
		Imports: make(map[string]string),
	}
}

func (this *Model) Template(tplfile string) string {
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

func (this *Model) Generate() {
	this.TableFields = this.GetTableFields(this.TableName)
	for _, tableField := range this.TableFields {
		tField := StructField{
			FieldName:         generator.CamelCase(tableField.Field),
			FieldType:         getFiledType(tableField),
			TableFieldName:    tableField.Field,
			TableFieldComment: tableField.Comment,
		}
		if tableField.Field == "delete_time" {
			this.HasDeleteTime = true
		}
		if tableField.Key == "PRI" {
			this.TablePrimaryId = tField.TableFieldName
			this.PrimaryIdType = tField.FieldType
			this.PrimaryId = generator.CamelCase(this.TablePrimaryId)
			tField.FieldName = this.PrimaryId
		}
		if tField.FieldType == "time.Time" {
			if _, ok := this.Imports[tField.FieldType]; !ok {
				this.Imports[tField.FieldType] = `"time"`
			}
		}
		this.StructFields = append(this.StructFields, tField)
	}
	if _, err := os.Stat(this.OutDir); os.IsNotExist(err) {
		if err = os.Mkdir(this.OutDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	tpl := this.Template("model.tpl")
	tmpl, err := template.New("model.tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	outfile := utils.File{Name: this.OutDir + utils.PathSeparator() + "model.go"}
	if outfile.Exist() {
		//生成最新ast buffer
		bufferbuf := new(bytes.Buffer)
		err = tmpl.Execute(bufferbuf, this)
		bufferfset := token.NewFileSet()
		bufferfd, err := parser.ParseFile(bufferfset, outfile.Name, bufferbuf.String(), 0)
		if err != nil {
			panic(err)
		}
		var bufferNode *ast.TypeSpec
		ast.Inspect(bufferfd, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				if x.Name.Name == this.StructName {
					bufferNode = x //找到struct node
				}
			}
			return true
		})
		//读取原文件
		oldfile, err := os.OpenFile(outfile.Name, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		oldsrc, err := ioutil.ReadAll(oldfile)
		if err != nil {
			panic(err)
		}
		oldfset := token.NewFileSet()
		oldfd, err := parser.ParseFile(oldfset, outfile.Name, oldsrc, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		//生成新文件
		newFileBuf := bytes.NewBufferString("")
		ast.Inspect(oldfd, func(n ast.Node) bool {
			switch oldNode := n.(type) {
			case *ast.File:
				newFileBuf = bytes.NewBufferString(`package ` + oldNode.Name.Name + "\n")
			case *ast.GenDecl:
				if oldNode.Tok == token.TYPE {
					for i, oldn := range oldNode.Specs {
						if oldn.(*ast.TypeSpec).Name.Name == this.StructName {
							oldNode.Specs[i] = bufferNode //将Buffer中的Node替换原Node
						}
					}
				}
				astToGo(newFileBuf, n)
			case *ast.FuncDecl:
				if oldNode.Doc != nil {
					length := len(oldNode.Doc.List)
					for i, doc := range oldNode.Doc.List {
						if i == 0 && 1 == length {
							newFileBuf.WriteString("\n" + doc.Text)
						} else if i == 0 {
							newFileBuf.WriteString("\n" + doc.Text + "\n")
						} else if i == (length - 1) {
							newFileBuf.WriteString(doc.Text)
						} else {
							newFileBuf.WriteString(doc.Text + "\n")
						}
					}
					oldNode.Doc = nil
				}
				astToGo(newFileBuf, oldNode)
			}
			return true
		})
		//ast.Print(oldfset, oldfd)
		//fmt.Println(newFileBuf)
		newFile, err := os.OpenFile(outfile.Name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		_, err = newFile.Write(newFileBuf.Bytes())
		if err != nil {
			panic(err)
		}
	} else {
		modelFile, err := os.OpenFile(outfile.Name, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		defer modelFile.Close()
		//生成model.go
		err = tmpl.Execute(modelFile, this)
		if err != nil {
			panic(err)
		}
	}
	outfile = utils.File{Name: this.OutDir + utils.PathSeparator() + "attributes.go"}
	attributesFile, err := os.OpenFile(outfile.Name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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
	fmt.Println("model: " + this.OutDir + " generate success!")
}

func GetTables(db *gorm.DB, database string, tableNames ...string) []Table {
	var (
		tables []Table
		d      *gorm.DB
	)
	if len(tableNames) == 0 {
		d = db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" + database + "';").Find(&tables)
	} else {
		d = db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE TABLE_NAME IN (" + strings.Join(tableNames, ",") + ") AND table_schema='" + database + "';").Find(&tables)
	}
	if d.Error != nil {
		panic(d.Error.Error())
	}
	return tables
}

//获取表所有字段信息
func (this *Model) GetTableFields(tableName string) []TableField {
	db := this.DB
	var fields []TableField
	d := db.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	if d.Error != nil {
		panic(d.Error.Error())
	}
	return fields
}

type Table struct {
	Name    string `gorm:"column:Name"`
	Comment string `gorm:"column:Comment"`
}

type StructField struct {
	FieldName         string
	FieldType         string
	TableFieldName    string
	TableFieldComment string
}

type TableField struct {
	Field      string `gorm:"column:Field"`
	Type       string `gorm:"column:Type"`
	Null       string `gorm:"column:Null"` //非空 YES/NO
	Key        string `gorm:"column:Key"`
	Default    string `gorm:"column:Default"`
	Extra      string `gorm:"column:Extra"`
	Privileges string `gorm:"column:Privileges"`
	Comment    string `gorm:"column:Comment"`
}

//获取字段类型
func getFiledType(field TableField) string {
	if field.Null == "YES" {
		return "interface{}"
	}
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
	case "binary":
		return "[]byte"
	default:
		return "string"
	}
}
