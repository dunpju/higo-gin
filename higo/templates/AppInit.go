package templates

import (
	"bytes"
	"github.com/dunpju/higo-utils/utils"
	"go/format"
	"go/token"
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
		getModule, err := utils.Mod.GetModule()
		if err != nil {
			panic(err)
		}
		moduleName = getModule
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
