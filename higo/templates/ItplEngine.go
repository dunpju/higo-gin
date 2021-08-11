package templates

const (
	controller  = "Controller"
	model       = "Model"
	attributes  = "attributes"
	NewFuncDecl = "func_decl"
)

type ItplEngine interface {
	Template(tplfile string) string
	Generate()
}
