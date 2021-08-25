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
	"strconv"
	"strings"
	"text/template"
)

type EnumMap struct {
	Key   string
	Value int
	Doc   string
}

func NewEnumMap(key string, value int, doc string) *EnumMap {
	return &EnumMap{Key: key, Value: value, Doc: doc}
}

type Enum struct {
	Package   string
	Name      string
	OutStruct string
	File      string
	Doc       string
	Enums     []*EnumMap
}

func NewEnum(pkg string, name string, file string) *Enum {
	reg := regexp.MustCompile(`(-e=[a-zA-Z]+\s*-f=).*`)
	if reg == nil {
		log.Fatalln("regexp err")
	}
	e := &Enum{}
	if fs := reg.FindString(name); fs != "" {
		name = strings.Trim(fs, " ")
		name = strings.Trim(name, "-e=")
		names := strings.Split(name, "-f=")
		if len(names) != 2 {
			log.Fatalln("name err")
		}
		name = strings.Trim(names[0], " ")
		docs := strings.Split(names[1], ":")
		doc := strings.Trim(docs[0], " ")
		e.Doc = doc
		es := strings.Split(docs[1], ",")
		for _, v := range es {
			em := strings.Split(v, "-")
			em1, err := strconv.Atoi(em[1])
			if err != nil {
				panic(err)
			}
			e.Enums = append(e.Enums, NewEnumMap(em[0], em1, em[2]))
		}
	}
	name = utils.Ucfirst(name)
	e.Name = name + enum
	e.OutStruct = file + utils.PathSeparator() + enum + strings.Trim(name, enum)
	e.File = e.OutStruct + ".go"
	e.Package = pkg
	log.Fatalln(e)
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
	outfile := utils.File{Name: this.File}
	if outfile.Exist() {
		log.Fatalln(this.File + " already existed")
	}
	outFile, err := os.OpenFile(this.File, os.O_RDWR|os.O_CREATE, 0755)
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
