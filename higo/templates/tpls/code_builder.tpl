package {{.Package}}

import (
	"gitee.com/dengpju/higo-code/code"
	. "github.com/dengpju/higo-gin/higo/errcode"
)

const (
	{{- range $i,$v := .KeyValueDocs}}
	{{- if eq $v.Iota "yes"}}
	{{$v.Key}} {{$.Name}} = iota + {{$v.Value}}  //{{$v.Doc}}
	{{- else}}
	{{- if eq $v.Value ""}}
	{{$v.Key}}  //{{$v.Doc}}
	{{- else}}
	{{$v.Key}} {{$.Name}} = {{$v.Value}}  //{{$v.Doc}}
	{{- end}}
	{{- end}}
    {{- end}}
)

func {{.FuncName}}() {
	code.Container().
    {{- range $i,$v := .KeyValueDocs}}
	    {{- if ne $i $.Len}}
	    Put({{$v.Key}}, "{{$v.Doc}}").
	    {{- else}}
	    Put({{$v.Key}}, "{{$v.Doc}}")
        {{- end}}
	{{- end}}
}