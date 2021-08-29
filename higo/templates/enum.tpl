package {{.Package}}

import "github.com/dengpju/higo-enum/enum"

//{{.Doc}}
type {{.Name}} int

func (this {{.Name}}) Name() string {
	return "{{.Name}}"
}

func (this {{.Name}}) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this {{.Name}}) Message() string {
	return enum.String(this)
}

const (
	{{- range $i,$v := .EnumMap}}
	{{$v.Key}} {{$.Name}} = {{$v.Value}} //{{$v.Doc}}
    {{- end}}
)

func (this {{$.Name}}) Register() enum.Message {
	return make(enum.Message).
	{{- range $i,$v := .EnumMap}}
	    {{- if ne $i $.LenMap}}
	    Put({{$v.Key}}, "{{$v.Doc}}").
	    {{- else}}
	    Put({{$v.Key}}, "{{$v.Doc}}")
        {{- end}}
	{{- end}}
}