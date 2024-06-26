package templates

import (
	"fmt"
	"github.com/dunpju/higo-gin/higo/templates/tpls"
	"github.com/dunpju/higo-orm/gen"
	"github.com/dunpju/higo-utils/utils"
	"github.com/dunpju/higo-utils/utils/fileutil"
	"github.com/dunpju/higo-utils/utils/stringutil"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"regexp"
	"strings"
)

type CodeMap struct {
	Key   string
	Value interface{}
	Doc   string
}

func NewCodeMap(key string, value interface{}, doc string) *CodeMap {
	return &CodeMap{Key: stringutil.Ucfirst(stringutil.CaseToCamel(key)), Value: value, Doc: doc}
}

type CodeBuilder struct {
	Package      string
	Name         string
	Len          int
	FuncName     string
	KeyValueDocs []KeyValueDoc
	File         string
	OutDir       string
	Arguments    *CodeArguments
}

func (this *CodeBuilder) generate() {
	keyValueDocs := this.KeyValueDocs
	this.Len = len(this.KeyValueDocs) - 1
	if len(keyValueDocs) > 0 {
		utils.Dir.Open(this.OutDir).Create()
		File := this.OutDir + utils.Dir.Separator() + this.File + ".go"
		this.FuncName = "code" + this.File
		if utils.File.Exist(File) && this.Arguments.Force != "yes" {
			log.Println(File + " already existed")
			return
		}
		outFile := utils.File.New(File, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
		defer outFile.Close()
		tmpl, err := this.Template("code_builder.tpl").Parse()
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

func (this *CodeBuilder) Template(tplfile string) *tpls.Tpl {
	return tpls.New(tplfile)
}

type Autoload struct {
	Package   string
	FuncNames []string
	OutDir    string
}

func NewAutoload(funcNames []string, out string) *Autoload {
	return &Autoload{Package: utils.Dir.Basename(out), FuncNames: funcNames, OutDir: out}
}

func (this *Autoload) generate() {
	utils.Dir.Open(this.OutDir).Create()
	File := this.OutDir + utils.Dir.Separator() + "Autoload.go"
	outFile := utils.File.New(File, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	defer outFile.Close()
	tmpl, err := this.Template("autoload.tpl").Parse()
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

func (this *Autoload) Template(tplfile string) *tpls.Tpl {
	return tpls.New(tplfile)
}

type CodeArguments struct {
	Package string
	Doc     string
	Name    string
	Out     string
	Auto    string
	Force   string
	Path    string
	Const   string
	Code    string
	Message string
	Iota    string
}

type Code struct {
	Package      string
	Name         string
	OutStruct    string
	OutDir       string
	File         string
	Doc          string
	RealName     string
	Code         string
	Iota         string
	LenMap       int
	CodeMap      []*CodeMap
	Codes        []*Code
	CodeBuilders []*CodeBuilder
	Arguments    CodeArguments
	funcNames    []string
}

// go run test\bin\main.go -gen=code -name=CodeErrorCode -out=test\app\Codes -auto=yes -force=yes -path=test\bin\yaml
// go run test\bin\main.go -gen=code -name=CodeErrorCode -out=test\app\Codes -auto=yes -force=yes -path=test\bin\yaml\400.yaml
// go run test\bin\main.go -gen=code -name=CodeErrorCode -out=test\app\Codes -auto=yes -force=yes -const=success -code=200 -message=成功 -iota=yes
func NewCode(args *CodeArguments) *Code {
	c := &Code{}
	c.Arguments = *args
	if args.Const != "" && args.Code != "" && args.Message != "" {
		c.Codes = append(c.Codes, newCode(args))
	} else if args.Path != "" {
		c.OutDir = args.Out
		inputFile := utils.File.Read(args.Path)
		if !inputFile.Exist() {
			log.Fatalln(args.Path + " configure file or directory non-exist")
		}
		if inputFile.IsDir() {
			files := utils.Dir.Open(args.Path).Suffix("yaml").Scan().Get()
			for _, file := range files {
				load(c, file, args)
			}
		} else {
			load(c, args.Path, args)
		}
	} else {
		log.Fatalln(" build parameter error")
	}
	return c
}

type KeyValueDoc struct {
	Key   interface{}
	Value interface{}
	Doc   string
	Iota  string
}

func load(c *Code, file string, args *CodeArguments) {
	fileName := utils.Dir.Basename(file, ".yaml")
	yamlFile := utils.File.Read(file)
	keySort := make([]string, 0)
	err := yamlFile.ForEach(func(line int, b []byte) bool {
		ok, err := regexp.Match("^[a-zA-Z]", b)
		if err != nil {
			panic(err)
		}
		if ok {
			keySort = append(keySort, strings.TrimRight(string(b), ":"))
		}
		return true
	})
	if err != nil {
		panic(err)
	}
	yamlFile = utils.File.Read(file)
	yamlMap := make(map[string]interface{})
	err = yaml.Unmarshal(yamlFile.Bytes(), yamlMap)
	if err != nil {
		log.Fatalln(err)
	}
	codeBuilder := &CodeBuilder{
		Package:   args.Package,
		Name:      args.Name,
		File:      fileName,
		OutDir:    args.Out,
		Arguments: args,
	}
	for i, k := range keySort {
		kk := utils.String.CaseToCamel(strings.ToLower(strings.Trim(k, "")))
		m := yamlMap[k]
		if m == nil {
			m = make((map[interface{}]interface{}))
		}
		v := m.(map[interface{}]interface{})
		iota, ok := v["iota"]
		if !ok {
			iota = "no"
		}
		code, ok := v["code"]
		if !ok {
			code = ""
		}
		if iota == "yes" && code == "" && i == 0 {
			ok, err := regexp.MatchString(`\d+`, fileName)
			if err != nil {
				panic(err)
			}
			if ok {
				code = fileName
			}
		}
		if code == "" {
			iota = "no"
		}
		message, ok := v["message"]
		if !ok {
			message = k + "错误"
		}
		codeBuilder.KeyValueDocs = append(codeBuilder.KeyValueDocs, KeyValueDoc{
			Key:   kk,
			Value: code,
			Doc:   message.(string),
			Iota:  iota.(string),
		})
	}
	c.CodeBuilders = append(c.CodeBuilders, codeBuilder)
}

func newCode(args *CodeArguments) *Code {
	c := &Code{
		Package:   args.Package,
		Doc:       args.Doc,
		Name:      args.Name,
		Code:      args.Code,
		Iota:      args.Iota,
		CodeMap:   make([]*CodeMap, 0),
		OutDir:    args.Out,
		OutStruct: args.Out + utils.Dir.Separator() + args.Name,
		RealName:  args.Name,
		File:      args.Out + utils.Dir.Separator() + args.Name + ".go",
		Arguments: *args,
	}
	c.CodeMap = append(c.CodeMap, NewCodeMap(args.Const, args.Code, args.Message))
	return c
}

func (this *Code) Template(tplfile string) *tpls.Tpl {
	return tpls.New(tplfile)
}

func (this *Code) Generate() {
	for _, e := range this.Codes {
		e.generate()
	}
	for _, b := range this.CodeBuilders {
		b.generate()
		this.funcNames = append(this.funcNames, b.FuncName)
	}
	if this.Arguments.Auto == gen.Yes && len(this.funcNames) > 0 {
		NewAutoload(this.funcNames, this.OutDir).generate()
	}
}

func (this *Code) generate() {
	utils.Dir.Open(this.OutDir).Create()
	if utils.File.Exist(this.File) && this.Arguments.Force != gen.Yes {
		log.Println(this.File + " already existed")
		return
	}
	outFile := fileutil.NewFile(this.File, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if !outFile.Exist() {
		outFile.Create()
	}
	defer outFile.Close()
	tmpl, err := this.Template("code.tpl").Parse()
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
