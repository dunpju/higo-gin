package templates

import "C"
import (
	"fmt"
	"github.com/dengpju/higo-gin/higo/templates/tpls"
	"github.com/dengpju/higo-utils/utils"
	"github.com/dengpju/higo-utils/utils/fileutil"
	"github.com/dengpju/higo-utils/utils/maputil"
	"github.com/dengpju/higo-utils/utils/stringutil"
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

type Code struct {
	Package   string
	Name      string
	OutStruct string
	OutDir    string
	File      string
	Doc       string
	RealName  string
	Code      string
	Iota      string
	LenMap    int
	CodeMap   []*CodeMap
	Codes     []*Code
	Arguments CodeArguments
}

type CodeBuilder struct {
	Package      string
	Name         string
	Len          int
	FuncName     string
	KeyValueDocs []maputil.KeyValueDoc
	OutDir       string
}

func NewCodeBuilder(name string, out string) *CodeBuilder {
	return &CodeBuilder{Package: utils.Dir.Basename(out), Name: name, OutDir: out}
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

var (
	codeRegexpStr = `(-c=[a-zA-Z_]+\s*-i=[0-9]+\s*-f=).*`
	funcNames     []string
)

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

// go run test\bin\main.go -gen=code -name="test\bin\yaml" -out=test\app\Codes -auto=yes -force=yes -extends=CodeErrorCode
// go run test\bin\main.go -gen=code -name=CodeErrorCode -out=test\app\Codes -auto=yes -force=yes -const=success -code=200 -message=成功
func NewCode(args *CodeArguments) *Code {
	c := &Code{}
	if args.Const != "" && args.Code != "" && args.Message != "" {
		c.Codes = append(c.Codes, newCode(args))
	} else if args.Path == "" {

	}

	//C := &Code{}
	//if fs := reg.FindString(args.Name); fs != "" {
	//	C.Codes = append(C.Codes, newCode(args))
	//} else {
	//	outfile := utils.File.Read(args.Name)
	//	if !outfile.Exist() {
	//		log.Fatalln(args.Name + " configure file non-exist")
	//	}
	//	if outfile.IsDir() {
	//		files := utils.Dir.Open(args.Name).Scan().Get()
	//		for _, filePath := range files {
	//			suffix := path.Ext(filePath)
	//			if ".md" == suffix {
	//				fileContext := string(utils.File.Read(filePath).ReadAll())
	//				build := NewCodeBuilder(args.Name, args.Out).parse(fileContext)
	//				build.generate()
	//				funcNames = append(funcNames, build.FuncName)
	//			} else if ".yaml" == suffix {
	//				//funcNames = append(funcNames, build.FuncName)
	//			}
	//		}
	//		if args.Auto == "yes" && len(funcNames) > 0 {
	//			NewAutoload(funcNames, args.Out).generate()
	//		}
	//	} else {
	//		err := outfile.ForEach(func(line int, b []byte) {
	//			s := string(b)
	//			s = strings.Replace(s, "\\", "", -1)
	//			s = strings.Trim(s, "\n")
	//			s = strings.Trim(s, "\r")
	//			s = strings.Trim(s, "")
	//			if "" != s {
	//				args.Name = s
	//				C.Codes = append(C.Codes, newCode(args))
	//			}
	//		})
	//		if err != nil {
	//			log.Fatalln(err)
	//		}
	//	}
	//}
	return c
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
}

func (this *Code) generate() {
	utils.Dir.Open(this.OutDir).Create()
	if utils.File.Exist(this.File) && this.Arguments.Force != "yes" {
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
