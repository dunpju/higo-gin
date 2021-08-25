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

type EnumMap struct {
	Key   string
	Value interface{}
	Doc   string
}

func NewEnumMap(key string, value interface{}, doc string) *EnumMap {
	return &EnumMap{Key: utils.Ucfirst(utils.CaseToCamel(key)), Value: value, Doc: doc}
}

type Enum struct {
	Package   string
	Name      string
	OutStruct string
	File      string
	Doc       string
	RealName  string
	EnumMap   []*EnumMap
	Enums     []*Enum
}

func NewEnum(pkg string, name string, file string) *Enum {
	reg := regexp.MustCompile(`(-e=[a-zA-Z_]+\s*-f=).*`)
	if reg == nil {
		log.Fatalln("regexp err")
	}
	e := &Enum{}
	if fs := reg.FindString(name); fs != "" {
		e.Enums = append(e.Enums, newEnum(pkg, name, file))
	} else {
		outfile := utils.NewFile(name)
		if !outfile.Exist() {
			log.Fatalln(name + " configure file non-exist")
		}
		outfile.ForEach(func(line int, s string) {
			s = strings.Replace(s, "\\", "", -1)
			s = strings.Trim(s, "\n")
			s = strings.Trim(s, "\r")
			s = strings.Trim(s, "")
			if "" != s {
				e.Enums = append(e.Enums, newEnum(pkg, s, file))
			}
		})
	}
	return e
}

func newEnum(pkg string, name string, file string) *Enum {
	reg := regexp.MustCompile(`(-e=[a-zA-Z_]+\s*-f=).*`)
	if reg == nil {
		log.Fatalln("regexp err")
	}
	name = strings.Replace(name, "\\", "", -1)
	name = strings.Trim(name, "\n")
	name = strings.Trim(name, "\r")
	name = strings.Trim(name, "")
	e := &Enum{}
	if fs := reg.FindString(name); fs != "" {
		name = strings.Trim(fs, " ")
		name = strings.Trim(name, "-e=")
		names := strings.Split(name, "-f=")
		if len(names) != 2 {
			log.Fatalln("name err")
		}
		name = strings.Trim(names[0], "")
		docs := strings.Split(names[1], ":")
		doc := strings.Trim(docs[0], "")
		e.Doc = doc
		es := strings.Split(docs[1], ",")
		for _, v := range es {
			em := strings.Split(v, "-")
			e.EnumMap = append(e.EnumMap,
				NewEnumMap(strings.Trim(em[0], ""),
					strings.Trim(em[1], ""),
					strings.Trim(strings.Trim(strings.Trim(em[2], "\n"), "\r"), " ")))
		}
		name = utils.Ucfirst(utils.CaseToCamel(name))
		e.Name = enum + name
		e.RealName = name
		e.OutStruct = file + utils.PathSeparator() + enum + strings.Trim(name, enum)
		e.File = e.OutStruct + ".go"
		e.Package = pkg
		return e
	} else {
		log.Fatalln(`name format error: ` + name)
	}
	return e
}

func (this *Enum) Template(tplfile string) string {
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

func (this *Enum) Generate() {
	for _, e := range this.Enums {
		e.generate()
	}
}

func (this *Enum) generate() {
	outFile, err := os.OpenFile(this.File, os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	tpl := this.Template("enum.tpl")
	tmpl, err := template.New("enum.tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	//生成enum
	err = tmpl.Execute(outFile, this)
	if err != nil {
		panic(err)
	}
	fmt.Println("enum: " + this.OutStruct + " generate success!")
}
