package templates

const (
	controller = "Controller"
	NewFuncDecl = "NewFuncDecl"
)

type ItplEngine interface {
	Template(tplfile string) string
	Generate()
}
