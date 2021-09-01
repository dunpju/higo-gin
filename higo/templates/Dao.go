package templates

import "C"
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
	Package   string
	Name      string
	OutStruct string
	OutDir    string
	File      string
	Doc       string
	RealName  string
	Iota      string
	LenMap    int
	DaoMap    []*DaoMap
	Daos      []*Dao
}

var daoRegexpStr = `(-c=[a-zA-Z_]+\s*-i=[0-9]+\s*-f=).*`

func NewDao(pkg string, name string, file string) *Dao {
	reg := regexp.MustCompile(daoRegexpStr)
	if reg == nil {
		log.Fatalln("regexp err")
	}
	D := &Dao{}
	if fs := reg.FindString(name); fs != "" {
		D.Daos = append(D.Daos, newDao(pkg, name, file))
	} else {
		outfile := utils.ReadFile(name)
		if !outfile.Exist() {
			log.Fatalln(name + " configure file non-exist")
		}
		outfile.ForEach(func(line int, s string) {
			s = strings.Replace(s, "\\", "", -1)
			s = strings.Trim(s, "\n")
			s = strings.Trim(s, "\r")
			s = strings.Trim(s, "")
			if "" != s {
				D.Daos = append(D.Daos, newDao(pkg, s, file))
			}
		})
	}
	return D
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
		name = strings.Trim(fs, "")
		name = strings.Trim(name, "-c=")
		codeName := strings.Split(name, "-i=")
		structName := codeName[0]
		flags := strings.Split(codeName[1], "-f=")
		if len(flags) != 2 {
			log.Fatalln("name err")
		}
		D.Iota = strings.Trim(flags[0], "")
		docs := strings.Split(flags[1], ":")
		doc := strings.Trim(docs[0], "")
		D.Doc = doc
		es := strings.Split(docs[1], ",")
		for _, v := range es {
			em := strings.Split(v, "-")
			k := strings.Trim(em[0], "")
			v := strings.Trim(D.Iota, "")
			d := strings.Trim(strings.Trim(strings.Trim(em[1], "\n"), "\r"), "")
			D.DaoMap = append(D.DaoMap, NewDaoMap(k, v, d))
		}
		D.LenMap = len(D.DaoMap) - 1
		name = utils.Ucfirst(utils.CaseToCamel(structName))
		D.Name = dao + name
		D.RealName = name
		D.OutDir = file
		D.OutStruct = D.OutDir + utils.PathSeparator() + dao + strings.Trim(name, dao)
		D.File = D.OutDir + utils.PathSeparator() + D.RealName + ".go"
		D.Package = pkg
		return D
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
	for _, e := range this.Daos {
		e.generate()
	}
}

func (this *Dao) generate() {
	utils.Dir(this.OutDir).Create()
	if utils.FileExist(this.File) {
		log.Println(this.File + " already existed")
		return
	}
	outFile := utils.NewFile(this.File, os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0755)
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
