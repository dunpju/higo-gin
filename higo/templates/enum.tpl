package {{.Package}}

import "github.com/dengpju/higo-enum/enum"

//{{.Doc}}
type {{.Name}} int

func (this {{.Name}}) Code() int64 {
	return enum.New(this).Code
}

func (this {{.Name}}) Message() string {
	return enum.New(this).Doc
}


const (
	{{- range $i,$v := .EnumMap}}
	{{- if eq 0 $i}}
	{{$v.Key}} {{$.Name}} = {{$v.Value}} //{{$v.Doc}}
	{{- else}}
	{{$v.Key}} = {{$v.Value}} //{{$v.Doc}}
	{{- end}}
    {{- end}}
)

func (this {{.Name}}) String() string {
	switch this {
	{{- range $v := .EnumMap}}
	case {{$v.Key}}:
    	return "{{$v.Doc}}"
    {{- end}}
	}
	return "未定义"
}