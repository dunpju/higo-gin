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
)

type VO struct {
	Package    string
	StructName string
	OutDir     string
	FileName   string
	JsonFile   string
	GormTag    bool
	JsonTag    bool
	Force      bool
	FieldList  []*VoField
}

func NewVO(name, out, jsonFile string, gormTag, jsonTag, force bool) *VO {
	unpreCamelName := utils.String.CaseToCamel(name)
	pkg := dirutil.Basename(out)
	outDir := out
	file := unpreCamelName + "VO.go"
	return &VO{
		Package:    pkg,
		StructName: unpreCamelName + "VO",
		OutDir:     outDir,
		FileName:   file,
		JsonFile:   jsonFile,
		GormTag:    gormTag,
		JsonTag:    jsonTag,
		Force:      force,
		FieldList:  make([]*VoField, 0)}
}

type VoField struct {
	FieldName string
	FieldType string
	Tag       string
}

func (this *VO) Template(tplfile string) *tpls.Tpl {
	return tpls.New(tplfile)
}

func (this *VO) Generate() {
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
		for fieldName, _ := range jsonContextMap {
			if fieldMaxLen < len(fieldName) {
				fieldMaxLen = len(fieldName)
			}
		}
		for fieldName, value := range jsonContextMap {
			var fieldType string
			switch value.(type) {
			case int:
				fieldType = "int"
			case int64, int32:
				fieldType = "int64"
			case uint64, uint32:
				fieldType = "uint64"
			case float64, float32:
				fieldType = "float64"
			case string:
				fieldType = "string"
			}
			camelFieldName := utils.String.CaseToCamel(fieldName)
			tag := ""
			if this.GormTag {
				tag += fmt.Sprintf(`gorm:"column:%s"`, fieldName)
			}
			if this.JsonTag {
				if tag != "" {
					tag += " "
					tag += LeftStrPad(fmt.Sprintf(`json:"%s"`, fieldName), fieldMaxLen-len(camelFieldName), " ")
				} else {
					tag += fmt.Sprintf(`json:"%s"`, fieldName)
				}
			}
			if tag != "" {
				tag = fmt.Sprintf("`%s`", tag)
			}
			this.FieldList = append(this.FieldList, &VoField{
				FieldName: camelFieldName + LeftStrPad(" ", fieldMaxLen-len(camelFieldName), " "),
				FieldType: fieldType,
				Tag:       tag,
			})
		}
	}

	if _, err := os.Stat(this.OutDir); os.IsNotExist(err) {
		if err = os.MkdirAll(this.OutDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	tmpl, err := this.Template("vo.tpl").Parse()
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
		fmt.Println("VO: " + outfile.Name + " was forced generate success!")
	} else {
		fmt.Println("VO: " + outfile.Name + " generate success!")
	}
}
