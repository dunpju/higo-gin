package templates

import "github.com/dengpju/higo-gin/higo/templates/tpls"

const (
	controller  = "Controller"
	model       = "Model"
	enum        = "Enum"
	code        = "Code"
	dao         = "Dao"
	attributes  = "attributes"
	NewFuncDecl = "func_decl"
)

type ItplEngine interface {
	Template(tplfile string) *tpls.Tpl
	Generate()
}
