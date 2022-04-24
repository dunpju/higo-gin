package templates

import (
	"bytes"
	"fmt"
	"github.com/dengpju/higo-gin/higo/templates/tpls"
	"github.com/dengpju/higo-utils/utils/dirutil"
	"github.com/dengpju/higo-utils/utils/fileutil"
	"github.com/dengpju/higo-utils/utils/stringutil"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Controller struct {
	Package   string
	Name      string
	OutStruct string
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
	return &Controller{Package: pkg, Name: name, OutStruct: outStruct, File: file}
}

func (this *Controller) Template(tplfile string) *tpls.Tpl {
	return tpls.New(tplfile)
}

func (this *Controller) Generate() {
	if fileutil.FileExist(this.File) {
		log.Fatalln(this.File + " already existed")
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
	_, mainfile, _, _ := runtime.Caller(4)
	app := strings.Trim(mainfile, "main.go") + ".." + dirutil.PathSeparator() + "app"
	beansGofile := app + dirutil.PathSeparator() + "Beans" + dirutil.PathSeparator() + "Bean.go"
	utifile := fileutil.File{Name: beansGofile}
	if !utifile.Exist() {
		log.Fatalln("Bean.go file non-existent, bean route cannot auto-load")
	}
	beansFile, err := os.OpenFile(beansGofile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	src, err := ioutil.ReadAll(beansFile)
	if err != nil {
		panic(err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}
	newPkgPath := GetModName() + "/" + strings.ReplaceAll(dirutil.Dirname(this.File), "\\", "/")
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