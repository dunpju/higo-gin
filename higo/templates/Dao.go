package templates

import (
	"fmt"
	"github.com/dengpju/higo-utils/utils"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
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
	PackageName       string
	Imports           map[string]string
	StructName        string
	ModelPackageName  string
	ModelName         string
	EntityPackageName string
	EntityName        string
	PrimaryId         string
	PrimaryIdType     string
	TablePrimaryId    string
	TableFields       []TableField
	ModelFields       []TplField
	HasDeleteTime     bool
	OutStruct         string
	OutDir            string
	File              string
}

var daoRegexpStr = `(-c=[a-zA-Z_]+\s*-i=[0-9]+\s*-f=).*`

func NewDao(pkg string, name string, file string) *Dao {

	return &Dao{}
}

func newDao(pkg string, name string, file string) *Dao {
	reg := regexp.MustCompile(daoRegexpStr)
	if reg == nil {
		log.Fatalln("regexp err")
	}
	name = strings.Replace(name, "\\", "", -1)
	name = strings.Trim(name, "\n")
	name = strings.Trim(name, "\r")
	name = strings.Trim(name, "")
	D := &Dao{}
	if fs := reg.FindString(name); fs != "" {

	} else {
		log.Fatalln(`name format error: ` + name)
	}
	return D
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
	if utils.FileExist(this.File) {
		log.Println(this.File + " already existed")
		return
	}
	outFile := utils.NewFile(this.File, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
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
	fmt.Println("dao: " + this.OutStruct + " generate success!")
}
