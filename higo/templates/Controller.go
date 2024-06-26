package templates

import (
	"bytes"
	"fmt"
	"github.com/dunpju/higo-gin/higo/templates/tpls"
	"github.com/dunpju/higo-orm/gen"
	"github.com/dunpju/higo-utils/utils"
	"github.com/dunpju/higo-utils/utils/dirutil"
	"github.com/dunpju/higo-utils/utils/fileutil"
	"github.com/dunpju/higo-utils/utils/stringutil"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type ControllerTool struct {
	Name                 string   // Controller Struct Name
	ParamTag             []string // List/Add/Edit/Delete
	ConfirmBeginGenerate gen.YesNo
	IsGenerateParam      gen.YesNo
	OutParamDir          string
}

func NewControllerTool() *ControllerTool {
	return &ControllerTool{ConfirmBeginGenerate: gen.Yes, IsGenerateParam: gen.Yes}
}

func (this *ControllerTool) Generate() {
	if this.IsGenerateParam.Bool() {
		for _, tag := range this.ParamTag {
			NewParam(this.Name+tag, this.OutParamDir, "", "", false).Generate()
		}
	}
}

type Controller struct {
	Package   string
	Name      string
	OutStruct string
	SelfName  string
	File      string
}

type FuncDecl struct {
	Recv     string
	FuncName string
	Results  string
	Returns  string
}

func NewController(pkg string, name string, file string) *Controller {
	name = stringutil.Ucfirst(name)
	name = name + controller
	outStruct := file + dirutil.PathSeparator() + strings.Replace(name, controller, "", -1) + controller
	file = outStruct + ".go"
	return &Controller{Package: pkg, Name: name, OutStruct: outStruct, SelfName: strings.ReplaceAll(outStruct, "\\", "/"), File: file}
}

func (this *Controller) Template(tplfile string) *tpls.Tpl {
	return tpls.New(tplfile)
}

func (this *Controller) Generate() {
	if fileutil.FileExist(this.File) {
		log.Println(this.File + " already existed")
		return
	}
	outFile := fileutil.NewFile(this.File, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if !outFile.Exist() {
		outFile.Create()
	}
	defer outFile.Close()
	tmpl, err := this.Template("controller.tpl").Parse()
	if err != nil {
		panic(err)
	}
	//生成controller
	err = tmpl.Execute(outFile, this)
	if err != nil {
		panic(err)
	}
	//bean route
	beansGofile := utils.Dir.Dirname(this.File) + dirutil.PathSeparator() + ".." + dirutil.PathSeparator() + "Beans" + dirutil.PathSeparator() + "Bean.go"
	utifile := fileutil.File{Name: beansGofile}
	if !utifile.Exist() {
		log.Fatalln("Bean.go file non-existent, bean route cannot auto-load")
	}
	beansFile, err := os.OpenFile(beansGofile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	src, err := io.ReadAll(beansFile)
	if err != nil {
		panic(err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	goMod := utils.Mod.GetModInfo()
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	childPath := utils.Mod.GetGoModChildPath(pwd)
	var childPathStr string
	if len(childPath) > 0 {
		childPathStr = fmt.Sprintf("/%s/", strings.Join(childPath, "/"))
	}
	afterPath := fmt.Sprintf("%s%s", childPathStr, strings.ReplaceAll(utils.Dir.Dirname(this.File), "\\", "/"))
	afterPathMatch, err := regexp.MatchString("^/", afterPath)
	if err != nil {
		panic(err)
	}
	if !afterPathMatch {
		afterPath = fmt.Sprintf("/%s", afterPath)
	}
	newPkgPath := goMod.Module.Path + afterPath
	funcName := strings.ReplaceAll(newPkgPath, "/", "_")
	funcName = strings.ReplaceAll(funcName, ".", "8")
	funcName = "New_gen_" + strings.ReplaceAll(funcName, "-", "9") + "_" + this.Name
	buffer := bytes.NewBufferString("")
	//import
	isImptHandle := false
	newImptSpec := &ast.ImportSpec{}
	recvTypeSpec := &ast.TypeSpec{}
	var (
		newFuncDeclOnce   sync.Once
		newFuncDeclFormat string
	)
	hasFuncDecl := false
	newFuncDecl := &FuncDecl{}
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			buffer = bytes.NewBufferString(`package ` + x.Name.Name + "\n")
		case *ast.GenDecl:
			if !isImptHandle {
				isImptHandle = true
				// 是否已经import
				hasImported := false
				for _, v := range x.Specs {
					if impt, ok := v.(*ast.ImportSpec); ok {
						// 如果已经包含包
						if impt.Path.Value == strconv.Quote(newPkgPath) {
							newImptSpec = impt
							hasImported = true
							break
						}
					}
				}
				if x.Tok == token.IMPORT {
					// 如果没有import，则import
					if !hasImported {
						newImptSpec = &ast.ImportSpec{
							Path: &ast.BasicLit{
								Kind:  token.STRING,
								Value: strconv.Quote(newPkgPath),
							},
						}
						x.Specs = append(x.Specs, newImptSpec)
					}
				}
			}
			if x.Tok == token.TYPE {
				recvTypeSpec = x.Specs[0].(*ast.TypeSpec)
			}
			astToGo(buffer, n)
		case *ast.FuncDecl:
			//判断是否存在，不能重复
			funcDeclBuf := bytes.NewBufferString("")
			err := format.Node(funcDeclBuf, token.NewFileSet(), n)
			if err != nil {
				panic(err)
			}
			if newImptSpec.Name != nil {
				newFuncDeclOnce.Do(func() {
					newFuncDeclFormat = fmt.Sprintf(funcDecl, recvTypeSpec.Name.String(),
						funcName, newImptSpec.Name.String()+".", this.Name)
					newFuncDecl.Recv = recvTypeSpec.Name.String()
					newFuncDecl.FuncName = funcName
					newFuncDecl.Results = newImptSpec.Name.String() + "." + this.Name
					newFuncDecl.Returns = newImptSpec.Name.String() + ".New" + this.Name + "()"
				})
			} else {
				newFuncDeclOnce.Do(func() {
					newFuncDeclFormat = fmt.Sprintf(funcDecl, recvTypeSpec.Name.String(),
						funcName, this.Package+".", this.Name)
					newFuncDecl.Recv = recvTypeSpec.Name.String()
					newFuncDecl.FuncName = funcName
					newFuncDecl.Results = this.Package + "." + this.Name
					newFuncDecl.Returns = this.Package + ".New" + this.Name + "()"
				})
			}
			if strings.Contains(funcDeclBuf.String(), newFuncDeclFormat) {
				hasFuncDecl = true
			}
			astToGo(buffer, n)
		}
		return true
	})
	//追加
	if !hasFuncDecl {
		tmpl, err := this.Template("func_decl.tpl").Parse()
		if err != nil {
			panic(err)
		}
		newFuncDeclBuffer := bytes.NewBufferString("")
		err = tmpl.Execute(newFuncDeclBuffer, newFuncDecl)
		if err != nil {
			panic(err)
		}
		newBuffer := bytes.NewBufferString(buffer.String() + "\n" + newFuncDeclBuffer.String())
		newBeansFile, err := os.OpenFile(beansGofile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		_, err = newBeansFile.Write(newBuffer.Bytes())
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("controller: " + this.OutStruct + " generate success!")
}
