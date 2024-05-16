package {{.Package}}

type {{.StructName}} struct {
	{{- range $i,$iter := .FieldList}}
	{{$iter.FieldName}}{{$iter.FieldType}} {{$iter.Tag}}
	{{- end}}
}