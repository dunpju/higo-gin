package templates

const (
	controller  = "Controller"
	model       = "Model"
	enum        = "Enum"
	attributes  = "attributes"
	NewFuncDecl = "func_decl"
)

type ItplEngine interface {
	Template(tplfile string) string
	Generate()
}
