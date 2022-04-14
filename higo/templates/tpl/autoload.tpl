package errcode

func autoload() {
    {{- range .FuncNames}}
	{{.}}()
    {{- end}}
}