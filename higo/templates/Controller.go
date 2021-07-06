package templates

import (
	"bytes"
	"fmt"
	"github.com/dengpju/higo-utils/utils"
	"github.com/pkg/errors"
	"github.com/win5do/go-lib/errx"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

type Controller struct {
	Package string
	Name    string
	File    string
}

type FuncDecl struct {
	Recv     string
	FuncName string
	Results  string
	Returns  string
}

func NewController(pkg string, name string, file string) *Controller {
	name = utils.Ucfirst(name)
	name = name + controller
	file = file + utils.PathSeparator() + strings.Trim(name, controller) + controller + ".go"
	return &Controller{Package: pkg, Name: name, File: file}
}

func (this *Controller) Template(tplfile string) string {
	_, file, _, _ := runtime.Caller(0)
	file = strings.TrimRight(file, ".go") + ".tpl"
	if tplfile == "NewFuncDecl.tpl" {
		file = path.Dir(file) + "/" + tplfile
	}
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

func (this *Controller) Generate() {
	controllerTpl, err := os.OpenFile(this.File, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer controllerTpl.Close()
	tpl := this.Template("")
	tmpl, err := template.New(controller).Parse(tpl)
	if err != nil {
		panic(err)
	}
	//生成controller
	err = tmpl.Execute(controllerTpl, this)
	if err != nil {
		panic(err)
	}
	//bean route
	_, mainfile, _, _ := runtime.Caller(3)
	app := strings.Trim(mainfile, "main.go") + ".." + utils.PathSeparator() + "app"
	beansGofile := app + utils.PathSeparator() + "Beans" + utils.PathSeparator() + "Bean.go"
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
	//log.Println("package name: ", this.Package)
	//log.Println("controller name: ", this.Name)
	//log.Println("output file: ", this.File)
	newPkgPath := GetModName() + "/" + strings.ReplaceAll(utils.Dirname(this.File), "\\", "/")
	buffer := bytes.NewBufferString("")
	//import
	isImptHandle := false
	newImptSpec := &ast.ImportSpec{}
	recvTypeSpec := &ast.TypeSpec{}
	var newFuncDeclOnce sync.Once
	newFuncDeclFormat := funcDecl
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
			_ = astToGo(buffer, n)
		case *ast.FuncDecl:
			//判断还是是否存在，不能重复
			funcDeclBuf := bytes.NewBufferString("")
			err := format.Node(funcDeclBuf, token.NewFileSet(), n)
			if err != nil {
				panic(err)
			}
			if newImptSpec.Name != nil {
				newFuncDeclOnce.Do(func() {
					newFuncDeclFormat = fmt.Sprintf(newFuncDeclFormat, recvTypeSpec.Name.String(),
						this.Name, newImptSpec.Name.String()+".", this.Name)
					newFuncDecl.Recv = recvTypeSpec.Name.String()
					newFuncDecl.FuncName = this.Name
					newFuncDecl.Results = newImptSpec.Name.String() + "." + this.Name
					newFuncDecl.Returns = newImptSpec.Name.String() + ".New" + this.Name + "()"
				})
			} else {
				newFuncDeclOnce.Do(func() {
					newFuncDeclFormat = fmt.Sprintf(newFuncDeclFormat, recvTypeSpec.Name.String(),
						this.Name, this.Package+".", this.Name)
					newFuncDecl.Recv = recvTypeSpec.Name.String()
					newFuncDecl.FuncName = this.Name
					newFuncDecl.Results = this.Package + "." + this.Name
					newFuncDecl.Returns = this.Package + ".New" + this.Name + "()"
				})
			}
			if strings.Contains(funcDeclBuf.String(), newFuncDeclFormat) {
				hasFuncDecl = true
			}
			_ = astToGo(buffer, n)
		}
		return true
	})
	//ast.Print(fset, f)
	//不存在则追加
	if !hasFuncDecl {
		tpl := this.Template("NewFuncDecl.tpl")
		tmpl, err := template.New(NewFuncDecl).Parse(tpl)
		if err != nil {
			panic(err)
		}
		newFuncDeclBuffer := bytes.NewBufferString("")
		//生成controller
		err = tmpl.Execute(newFuncDeclBuffer, newFuncDecl)
		if err != nil {
			panic(err)
		}
		newBuffer := bytes.NewBufferString(buffer.String() + "\n" + newFuncDeclBuffer.String())
		fmt.Println(newBuffer.String())
		newBeansFile, err := os.OpenFile(beansGofile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		n, err := newBeansFile.Write(newBuffer.Bytes())
		if err != nil {
			panic(err)
		}
		fmt.Println(n)
	}
	//fmt.Println(os.Stdout.Write(buffer.Bytes()))
}

func run(dir string, pkgName string) error {
	/**
	dir, err := getImportPkg("go.uber.org/zap")
	if err != nil {
		return errx.WithStackOnce(err)
	}
	log.Printf("dir: %+v", dir)
	*/

	pkg, err := parseDir(dir, pkgName)
	if err != nil {
		return errx.WithStackOnce(err)
	}
	//funcs, err := walkAst(pkg)
	_, err = walkAst(pkg)
	if err != nil {
		return errx.WithStackOnce(err)
	}
	/**
	err = writeGoFile(os.Stdout, funcs)
	if err != nil {
		return errx.WithStackOnce(err)
	}
	*/
	return nil
}

func getImportPkg(pkg string) (string, error) {
	p, err := build.Import(pkg, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, err
}

func parseDir(dir, pkgName string) (*ast.Package, error) {
	pkgMap, err := parser.ParseDir(
		token.NewFileSet(),
		dir,
		func(info os.FileInfo) bool {
			// skip go-test
			return !strings.Contains(info.Name(), "_test.go")
		},
		parser.Mode(0), // no comment
	)
	if err != nil {
		return nil, errx.WithStackOnce(err)
	}
	pkg, ok := pkgMap[pkgName]
	if !ok {
		err := errors.New("not found")
		return nil, errx.WithStackOnce(err)
	}
	return pkg, nil
}

type visitor struct {
	funcs []*ast.FuncDecl
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Recv == nil ||
			!n.Name.IsExported() ||
			len(n.Recv.List) != 1 {
			return nil
		}
		t, ok := n.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return nil
		}
		fmt.Println(t.X.(*ast.Ident).String())
		/**
		if t.X.(*ast.Ident).String() != "SugaredLogger" {
			return nil
		}
		*/
		//log.Printf("func name: %s", n.Name.String())
		v.funcs = append(v.funcs, rewriteFunc(n))
	}
	return v
}

func walkAst(node ast.Node) ([]ast.Decl, error) {
	v := &visitor{}
	ast.Walk(v, node)
	//log.Printf("funcs len: %d", len(v.funcs))
	var decls []ast.Decl
	for _, v := range v.funcs {
		decls = append(decls, v)
	}
	return decls, nil
}

func rewriteFunc(fn *ast.FuncDecl) *ast.FuncDecl {
	fn.Recv = nil
	fnName := fn.Name.String()
	var args []string
	for _, field := range fn.Type.Params.List {
		for _, id := range field.Names {
			idStr := id.String()
			_, ok := field.Type.(*ast.Ellipsis)
			if ok {
				// Ellipsis args
				idStr += "..."
			}
			args = append(args, idStr)
		}
	}
	exprStr := fmt.Sprintf(`_globalS.%s(%s)`, fnName, strings.Join(args, ","))
	expr, err := parser.ParseExpr(exprStr)
	if err != nil {
		panic(err)
	}
	var body []ast.Stmt
	if fn.Type.Results != nil {
		body = []ast.Stmt{
			&ast.ReturnStmt{
				// Return:
				Results: []ast.Expr{expr},
			},
		}
	} else {
		body = []ast.Stmt{
			&ast.ExprStmt{
				X: expr,
			},
		}
	}
	fn.Body.List = body
	return fn
}

// Output Go code
func writeGoFile(wr io.Writer, funcs []ast.Decl) error {
	header := `// Code generated by log-gen. DO NOT EDIT.
package logx
`
	buffer := bytes.NewBufferString(header)
	for _, fn := range funcs {
		err := astToGo(buffer, fn)
		if err != nil {
			return errx.WithStackOnce(err)
		}
	}
	_, err := wr.Write(buffer.Bytes())
	return err
}

func astToGo(dst *bytes.Buffer, node interface{}) error {
	addNewline := func() {
		err := dst.WriteByte('\n') // add newline
		if err != nil {
			log.Panicln(err)
		}
	}
	addNewline()
	err := format.Node(dst, token.NewFileSet(), node)
	if err != nil {
		return err
	}
	addNewline()
	return nil
}
