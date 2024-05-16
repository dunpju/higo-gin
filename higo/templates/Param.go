package templates

import (
	"encoding/json"
	"fmt"
	"github.com/dunpju/higo-gin/higo/templates/tpls"
	"github.com/dunpju/higo-utils/utils"
	"github.com/dunpju/higo-utils/utils/dirutil"
	"github.com/dunpju/higo-utils/utils/fileutil"
	"log"
	"os"
	"sort"
)

const (
	List   = "List"
	Add    = "Add"
	Edit   = "Edit"
	Delete = "Delete"
)

type Param struct {
	Package              string
	StructName           string
	LowerCamelStructName string
	OutDir               string
	FileName             string
	JsonFile             string
	Tag                  string
	Force                bool
	ParamFieldList       []*ParamField
}

func NewParam(name, out, jsonFile, tag string, force bool) *Param {
	unpreCamelName := utils.String.CaseToCamel(name)
	pkg := dirutil.Basename(out)
	outDir := out
	file := "Param" + unpreCamelName + ".go"
	return &Param{Package: pkg,
		StructName:           unpreCamelName,
		LowerCamelStructName: utils.String.Lcfirst(unpreCamelName),
		OutDir:               outDir,
		FileName:             file,
		JsonFile:             jsonFile,
		Tag:                  tag,
		Force:                force,
		ParamFieldList:       make([]*ParamField, 0)}
}

func (this *Param) Template(tplfile string) *tpls.Tpl {
	return tpls.New(tplfile)
}

type ParamField struct {
	FieldName string
	FieldType string
	Tag       string
	TagName   string
}

func (this *Param) Generate() {
	outFile := this.OutDir + dirutil.PathSeparator() + this.FileName
	if fileutil.FileExist(outFile) && !this.Force {
		log.Println(outFile + " already existed")
		return
	}

	if this.JsonFile != "" {
		if !fileutil.FileExist(this.JsonFile) {
			log.Println(this.JsonFile + " nonexistent")
			return
		}
		jsonContextMap := make(map[string]interface{})
		jsonContext := utils.File.Read(this.JsonFile)
		err := json.Unmarshal(jsonContext.Bytes(), &jsonContextMap)
		if err != nil {
			panic(err)
		}
		fieldMaxLen := 0
		fieldNameSort := make([]string, 0)
		for fieldName, _ := range jsonContextMap {
			fieldNameSort = append(fieldNameSort, fieldName)
			if fieldMaxLen < len(fieldName) {
				fieldMaxLen = len(fieldName)
			}
		}
		sort.Strings(fieldNameSort)
		for _, fieldName := range fieldNameSort {
			value := jsonContextMap[fieldName]
			fieldType := TypeAssert(value)
			tag := "json"
			if this.Tag != "" {
				tag = this.Tag
			}
			camelFieldName := utils.String.CaseToCamel(fieldName)
			this.ParamFieldList = append(this.ParamFieldList, &ParamField{
				FieldName: camelFieldName + LeftStrPad(" ", fieldMaxLen-len(camelFieldName), " "),
				FieldType: fieldType,
				Tag:       tag,
				TagName:   fieldName,
			})
		}
	}

	if _, err := os.Stat(this.OutDir); os.IsNotExist(err) {
		if err = os.MkdirAll(this.OutDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	tmpl, err := this.Template("param.tpl").Parse()
	if err != nil {
		panic(err)
	}
	outfile := fileutil.File{Name: this.OutDir + dirutil.PathSeparator() + this.FileName}
	paramFile := fileutil.NewFile(outfile.Name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if !paramFile.Exist() {
		paramFile.Create()
	}
	defer paramFile.Close()
	//生成param.go
	err = tmpl.Execute(paramFile, this)
	if err != nil {
		panic(err)
	}
	if this.Force {
		fmt.Println("param: " + outfile.Name + " was forced generate success!")
	} else {
		fmt.Println("param: " + outfile.Name + " generate success!")
	}
}

// LeftStrPad
// input string 原字符串
// padLength int 规定补齐后的字符串位数
// padString string 自定义填充字符串
func LeftStrPad(input string, padLength int, padString string) string {
	output := ""
	for i := 1; i <= padLength; i++ {
		output += padString
	}
	return output + input
}
