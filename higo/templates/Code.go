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

type CodeMap struct {
	Key   string
	Value interface{}
	Doc   string
}

func NewCodeMap(key string, value interface{}, doc string) *CodeMap {
	return &CodeMap{Key: utils.Ucfirst(utils.CaseToCamel(key)), Value: value, Doc: doc}
}

type Code struct {
	Package   string
	Name      string
	OutStruct string
	OutDir    string
	File      string
	Doc       string
	RealName  string
	CodeType  string
	LenMap    int
	CodeMap   []*CodeMap
	Codes     []*Code
}

func NewCode(pkg string, name string, file string) *Code {
	reg := regexp.MustCompile(`(-e=[a-zA-Z_]+\s*-f=).*`)
	if reg == nil {
		log.Fatalln("regexp err")
	}
	e := &Code{}
	if fs := reg.FindString(name); fs != "" {
		e.Codes = append(e.Codes, newCode(pkg, name, file))
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
				e.Codes = append(e.Codes, newCode(pkg, s, file))
			}
		})
	}
	return e
}

func newCode(pkg string, name string, file string) *Code {
	reg := regexp.MustCompile(`(-e=[a-zA-Z_]+\s*-f=).*`)
	if reg == nil {
		log.Fatalln("regexp err")
	}
	name = strings.Replace(name, "\\", "", -1)
	name = strings.Trim(name, "\n")
	name = strings.Trim(name, "\r")
	name = strings.Trim(name, "")
	E := &Code{}
	if fs := reg.FindString(name); fs != "" {
		name = strings.Trim(fs, "")
		name = strings.Trim(name, "-e=")
		names := strings.Split(name, "-f=")
		if len(names) != 2 {
			log.Fatalln("name err")
		}
		name = strings.Trim(names[0], "")
		docs := strings.Split(names[1], ":")
		doc := strings.Trim(docs[0], "")
		E.Doc = doc
		characterReg := regexp.MustCompile(`([a-zA-Z_]).*`)
		if characterReg == nil {
			log.Fatalln("character regexp err")
		}
		es := strings.Split(docs[1], ",")
		for _, v := range es {
			em := strings.Split(v, "-")
			k := strings.Trim(em[0], "")
			v := strings.Trim(em[1], "")
			d := strings.Trim(strings.Trim(strings.Trim(em[2], "\n"), "\r"), "")
			E.CodeMap = append(E.CodeMap, NewCodeMap(k, v, d))
			if valueMatch := characterReg.FindString(v); valueMatch != "" {
				E.CodeType = "string"
			} else {
				E.CodeType = "int"
			}
		}
		E.LenMap = len(E.CodeMap) - 1
		name = utils.Ucfirst(utils.CaseToCamel(name))
		E.Name = code + name
		E.RealName = name
		E.OutDir = file + utils.PathSeparator() + code + E.RealName
		E.OutStruct = E.OutDir + utils.PathSeparator() + code + strings.Trim(name, code)
		E.File = E.OutDir + utils.PathSeparator() + "code.go"//TODO
		E.Package = code + name
		return E
	} else {
		log.Fatalln(`name format error: ` + name)
	}
	return E
}

func (this *Code) Template(tplfile string) string {
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

func (this *Code) Generate() {
	for _, e := range this.Codes {
		e.generate()
	}
}

func (this *Code) generate() {
	utils.Dir(this.OutDir).Create()
	utils.FileFlag = os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	utils.SetModePerm(0755)
	outfile := utils.File{Name: this.File}
	if outfile.Exist() {
		log.Println(this.File + " already existed")
		return
	}
	outFile := utils.NewFile(this.File)
	defer outFile.Close()
	tpl := this.Template("code.tpl")
	tmpl, err := template.New("code.tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	//生成enum
	err = tmpl.Execute(outFile.File(), this)
	if err != nil {
		panic(err)
	}
	fmt.Println("code: " + this.OutStruct + " generate success!")
}
