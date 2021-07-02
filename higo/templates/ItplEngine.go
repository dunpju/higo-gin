package templates

const (
	controller = "Controller"
)

type ItplEngine interface {
	Template() string
	Generate()
}
