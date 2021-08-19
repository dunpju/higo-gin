package {{.Package}}

import "github.com/dengpju/higo-enum/enum"

type {{.Name}} int

func (this {{.Name}}) Code() int64 {
	return enum.New(this).Code
}

func (this {{.Name}}) Message() string {
	return enum.New(this).Doc
}


const (
	{{- range $i, $v := .Enums}}
	{{if eq 0 $i}}
	{{$v.Key}} = iota + {{.InitIncr}}
	{{else}}
	{{$v}}
    {{end}}
    {{- end}}
)

func (this {{.Name}}) String() string {
	switch this {
	{{- range _, $v := .Enums}}
	case {{$v.Key}}:
    		return "{{$v.Value}}"
    {{- end}}
	}
	return "未定义"
}