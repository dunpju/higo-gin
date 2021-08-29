package templates

const (
	controller  = "Controller"
	model       = "Model"
	enum        = "Enum"
	code        = "Code"
	attributes  = "attributes"
	NewFuncDecl = "func_decl"
)

type ItplEngine interface {
	Template(tplfile string) string
	Generate()
}
