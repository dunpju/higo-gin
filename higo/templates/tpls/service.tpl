package {{.Package}}

import (
	"github.com/dengpju/higo-gin/higo"
)

const Self{{.Name}} = "{{.SelfName}}"

type {{.Name}} struct {
}

func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{}
}

func (this *{{.Name}}) New() higo.IClass {
	return New{{.Name}}()
}