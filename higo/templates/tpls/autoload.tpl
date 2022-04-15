package {{.Package}}

func autoload() {
    {{- range .FuncNames}}
	{{.}}()
    {{- end}}
}