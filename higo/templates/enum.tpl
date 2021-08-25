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
	{{- range $v := .EnumMap}}
	{{$.RealName}}{{$v.Key}} = {{$v.Value}} //{{$v.Doc}}
    {{- end}}
)

func (this {{.Name}}) String() string {
	switch this {
	{{- range $v := .EnumMap}}
	case {{$.RealName}}{{$v.Key}}:
    		return "{{$v.Doc}}"
    {{- end}}
	}
	return "未定义"
}