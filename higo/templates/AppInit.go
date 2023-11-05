package templates

import (
	"bytes"
	"go/format"
	"go/token"
	"runtime/debug"
	"sync"
)

const (
	module   = "module "
	funcDecl = "func (this *%s) %s() *%s%s"
)

var (
	moduleName = ""
	moduleOnce sync.Once
)

// GetModName 获取模块名称
func GetModName() string {
	moduleOnce.Do(func() {
		info, _ := debug.ReadBuildInfo()
		for _, dep := range info.Deps {
			if dep.Version == "(devel)" {
				moduleName = dep.Path
			}
		}
	})
	return moduleName
}

func astToGo(dst *bytes.Buffer, node interface{}) {
	addNewline := func() {
		err := dst.WriteByte('\n') // add newline
		if err != nil {
			panic(err)
		}
	}
	addNewline()
	err := format.Node(dst, token.NewFileSet(), node)
	if err != nil {
		panic(err)
	}
	addNewline()
}
