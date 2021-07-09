package templates

const (
	controller  = "controller"
	model       = "model"
	attributes  = "attributes"
	NewFuncDecl = "func_decl"
)

type ItplEngine interface {
	Template(tplfile string) string
	Generate()
}
