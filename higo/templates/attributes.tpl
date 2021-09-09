package {{.Package}}

import (
	"github.com/dengpju/higo-gin/higo"
	{{- range $impo := .Imports}}
    {{$impo}}
    {{- end}}
)

{{- range .TplFields}}
func With{{.Field}}(v {{.Type}}) higo.Property {
	return func(class higo.IClass) {
		class.(*{{$.StructName}}).{{.Field}} = v
	}
}
{{end}}
