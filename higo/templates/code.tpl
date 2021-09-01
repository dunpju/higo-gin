package {{.Package}}

import "gitee.com/dengpju/higo-code/code"

//{{.Doc}}
type {{.Name}} int64

func (this {{.Name}}) Message() string {
	return code.Get(this)
}

const (
	{{- range $i,$v := .CodeMap}}
	{{- if eq $i 0}}
	{{$v.Key}} {{$.Name}} = iota + {{$.Iota}}  //{{$v.Doc}}
	{{- else}}
	{{$v.Key}}  //{{$v.Doc}}
	{{- end}}
    {{- end}}
)

func (this {{$.Name}}) Register() code.Message {
	return code.Container().
	{{- range $i,$v := .CodeMap}}
	    {{- if ne $i $.LenMap}}
	    Put({{$v.Key}}, "{{$v.Doc}}").
	    {{- else}}
	    Put({{$v.Key}}, "{{$v.Doc}}")
        {{- end}}
	{{- end}}
}