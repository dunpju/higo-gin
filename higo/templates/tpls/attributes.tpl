package {{.PackageName}}

import (
	"github.com/dunpju/higo-gin/higo"
	{{- range $impo := .Imports}}
    {{$impo}}
    {{- end}}
)

const (
    {{- range .StructFields}}
    {{.FieldName}} = "{{.TableFieldName}}"  //{{.TableFieldComment}}
    {{- end}}
)


{{- range .StructFields}}
func With{{.FieldName}}(v {{.FieldType}}) higo.Property {
	return func(class higo.IClass) {
		class.(*{{$.StructName}}).{{.FieldName}} = v
	}
}
{{end}}
