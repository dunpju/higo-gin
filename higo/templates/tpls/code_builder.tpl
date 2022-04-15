package {{.Package}}

import "gitee.com/dengpju/higo-code/code"

const (
	{{- range $i,$v := .KeyValueDocs}}
	{{- if eq $i 0}}
	{{$v.Key}} {{$.Name}} = iota + {{$v.Value}}  //{{$v.Doc}}
	{{- else}}
	{{$v.Key}}  //{{$v.Doc}}
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