package templates

import (
	"fmt"
	"github.com/dengpju/higo-utils/utils"
	"github.com/dengpju/higo-utils/utils/dirutil"
	"github.com/dengpju/higo-utils/utils/fileutil"
	"github.com/dengpju/higo-utils/utils/maputil"
	"github.com/dengpju/higo-utils/utils/stringutil"
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
	return &CodeMap{Key: stringutil.Ucfirst(stringutil.CaseToCamel(key)), Value: value, Doc: doc}
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

type CodeBuilder struct {
	Name         string
	Len          int
	FuncName     string
	KeyValueDocs []maputil.KeyValueDoc
	OutDir       string
}

func NewCodeBuilder(name string) *CodeBuilder {
	return &CodeBuilder{Name: name}
}

func (this *CodeBuilder) parse(doc string) *CodeBuilder {
	messagePattern := `(@Message\(").*?("\))`
	messageRegexp := regexp.MustCompile(messagePattern)
	if messageRegexp == nil {
		log.Fatalln("regexp err")
	}
	messages := messageRegexp.FindAllStringSubmatch(doc, -1)
	message := make([]string, 0)
	for _, msg := range messages {
		msg0 := strings.Replace(msg[0], "\n", "", -1)
		msg0 = strings.Replace(msg0, `@Message("`, "", -1)
		msg0 = strings.Replace(msg0, `")`, "", -1)
		message = append(message, msg0)
	}

	codePattern := `(const\s+).*?(;)`
	codeRegexp := regexp.MustCompile(codePattern)
	if codeRegexp == nil {
		log.Fatalln("regexp err")
	}
	codes := codeRegexp.FindAllStringSubmatch(doc, -1)
	codd := make([]maputil.KeyValueDoc, 0)
	for i, cod := range codes {
		cod0 := strings.Replace(cod[0], "\n", "", -1)
		cod0 = strings.Replace(cod0, `const `, "", -1)
		cod0 = strings.Replace(cod0, `;`, "", -1)
		codSplit := strings.Split(cod0, "=")
		key := utils.String.CaseToCamel(strings.ToLower(strings.Trim(codSplit[0], "")))
		value := strings.ToLower(strings.Trim(codSplit[1], ""))
		codd = append(codd, *utils.Map.NewKeyValueDoc(key, value, message[i]))
	}
	this.KeyValueDocs = codd
	return this
}

func (this *CodeBuilder) generate() {
	keyValueDocs := this.KeyValueDocs
	this.Len = len(this.KeyValueDocs) - 1
	if len(keyValueDocs) > 0 {
		utils.Dir.Open(this.OutDir).Create()
		File := this.OutDir + utils.Dir.Separator() + strings.Replace(keyValueDocs[0].Value.(string), " ", "", -1) + ".go"
		this.FuncName = "code" + strings.Replace(keyValueDocs[0].Value.(string), " ", "", -1)
		if utils.File.Exist(File) {
			log.Println(File + " already existed")
			return
		}
		outFile := utils.File.New(File, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
		defer outFile.Close()
		tpl := this.Template("code.tpl")
		tmpl, err := template.New("code.tpl").Parse(tpl)
		if err != nil {
			panic(err)
		}
		//生成
		err = tmpl.Execute(outFile.File(), this)
		if err != nil {
			panic(err)
		}
		fmt.Println("code file " + File + " generate success!")
	}
}

func (this *CodeBuilder) Template(tplfile string) string {
	return tpl(tplfile)
}

func tpl(tplfile string) string {
	_, file, _, _ := runtime.Caller(0)
	file = path.Dir(file) + utils.Dir.Separator() + tplfile
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

type Autoload struct {
	FuncNames []string
	OutDir    string
}

func NewAutoload(funcNames []string) *Autoload {
	return &Autoload{FuncNames: funcNames}
}

func (this *Autoload) generate() {
	utils.Dir.Open(this.OutDir).Create()
	File := this.OutDir + utils.Dir.Separator() + "Autoload.go"
	outFile := utils.File.New(File, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	defer outFile.Close()
	tpl := this.Template("autoload.tpl")
	tmpl, err := template.New("autoload.tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	//生成
	err = tmpl.Execute(outFile.File(), this)
	if err != nil {
		panic(err)
	}
	fmt.Println("Autoload file " + File + " generate success!")
}

func (this *Autoload) Template(tplfile string) string {
	return tpl(tplfile)
}

var (
	codeRegexpStr = `(-c=[a-zA-Z_]+\s*-i=[0-9]+\s*-f=).*`
	autoRegexpStr = `(\s*-auto).*`
	funcNames     []string
)

// go run test\bin\main.go -gen=code -name="test\bin\md -auto" -out=test\app\Codes
func NewCode(pkg string, name string, file string) *Code {
	reg := regexp.MustCompile(codeRegexpStr)
	if reg == nil {
		log.Fatalln("code regexp err")
	}
	C := &Code{}
	if fs := reg.FindString(name); fs != "" {
		C.Codes = append(C.Codes, newCode(pkg, name, file))
	} else {
		reg := regexp.MustCompile(autoRegexpStr)
		if reg == nil {
			log.Fatalln("auto regexp err")
		}
		if auto := reg.ReplaceAllString(name, ""); auto != "" {
			name = auto
		}
		outfile := utils.File.Read(name)
		if !outfile.Exist() {
			log.Fatalln(name + " configure file non-exist")
		}
		if outfile.IsDir() {
			files := utils.Dir.Open(name).Suffix("md").Scan().Get()
			for _, filePath := range files {
				fileContext := string(utils.File.Read(filePath).ReadAll())
				build := NewCodeBuilder("CodeErrorCode").parse(fileContext)
				build.generate()
				funcNames = append(funcNames, build.FuncName)
			}
			log.Fatalln("ddd")
			if len(funcNames) > 0 {
				NewAutoload(funcNames).generate()
			}
		} else {
			err := outfile.ForEach(func(line int, b []byte) {
				s := string(b)
				s = strings.Replace(s, "\\", "", -1)
				s = strings.Trim(s, "\n")
				s = strings.Trim(s, "\r")
				s = strings.Trim(s, "")
				if "" != s {
					C.Codes = append(C.Codes, newCode(pkg, s, file))
				}
			})
			if err != nil {
				log.Fatalln(err)
			}
		}
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
		name = stringutil.Ucfirst(stringutil.CaseToCamel(structName))
		C.Name = code + name
		C.RealName = name
		C.OutDir = file
		C.OutStruct = C.OutDir + dirutil.PathSeparator() + code + strings.Trim(name, code)
		C.File = C.OutDir + dirutil.PathSeparator() + C.RealName + ".go"
		C.Package = pkg
		return C
	} else {
		log.Fatalln(`name format error: ` + name)
	}
	return C
}

func (this *Code) Template(tplfile string) string {
	_, file, _, _ := runtime.Caller(0)
	file = path.Dir(file) + dirutil.PathSeparator() + tplfile
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
	dirutil.Dir(this.OutDir).Create()
	if fileutil.FileExist(this.File) {
		log.Println(this.File + " already existed")
		return
	}
	outFile := fileutil.NewFile(this.File, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if !outFile.Exist() {
		outFile.Create()
	}
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
