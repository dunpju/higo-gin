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
	Iota      string
	LenMap    int
	CodeMap   []*CodeMap
	Codes     []*Code
}

var codeRegexpStr = `(-c=[a-zA-Z_]+\s*-i=[0-9]+\s*-f=).*`

func NewCode(pkg string, name string, file string) *Code {
	reg := regexp.MustCompile(codeRegexpStr)
	if reg == nil {
		log.Fatalln("regexp err")
	}
	C := &Code{}
	if fs := reg.FindString(name); fs != "" {
		C.Codes = append(C.Codes, newCode(pkg, name, file))
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
				C.Codes = append(C.Codes, newCode(pkg, s, file))
			}
		})
	}
	return C
}

func newCode(pkg string, name string, file string) *Code {
	reg := regexp.MustCompile(codeRegexpStr)
	if reg == nil {
		log.Fatalln("regexp err")
	}
	name = strings.Replace(name, "\\", "", -1)
	name = strings.Trim(name, "\n")
	name = strings.Trim(name, "\r")
	name = strings.Trim(name, "")
	C := &Code{}
	if fs := reg.FindString(name); fs != "" {
		name = strings.Trim(fs, "")
		name = strings.Trim(name, "-c=")
		codeName := strings.Split(name, "-i=")
		structName := codeName[0]
		flags := strings.Split(codeName[1], "-f=")
		if len(flags) != 2 {
			log.Fatalln("name err")
		}
		C.Iota = strings.Trim(flags[0], "")
		docs := strings.Split(flags[1], ":")
		doc := strings.Trim(docs[0], "")
		C.Doc = doc
		es := strings.Split(docs[1], ",")
		for _, v := range es {
			em := strings.Split(v, "-")
			k := strings.Trim(em[0], "")
			v := strings.Trim(C.Iota, "")
			d := strings.Trim(strings.Trim(strings.Trim(em[1], "\n"), "\r"), "")
			C.CodeMap = append(C.CodeMap, NewCodeMap(k, v, d))
		}
		C.LenMap = len(C.CodeMap) - 1
		name = utils.Ucfirst(utils.CaseToCamel(structName))
		C.Name = code + name
		C.RealName = name
		C.OutDir = file
		C.OutStruct = C.OutDir + utils.PathSeparator() + code + strings.Trim(name, code)
		C.File = C.OutDir + utils.PathSeparator() + C.RealName + ".go"
		C.Package = code + name
		return C
	} else {
		log.Fatalln(`name format error: ` + name)
	}
	return C
}

func (this *Code) Template(tplfile string) string {
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
