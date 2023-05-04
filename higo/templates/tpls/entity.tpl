package {{.PackageName}}

import (
	"github.com/dunpju/higo-gin/higo"
	{{- range $impo := .Imports}}
    {{$impo}}
    {{- end}}
)

type {{.StructName}} struct {
	isEdit      bool
	currentFlag higo.Flag
	{{- range .StructFields}}
	{{.FieldName}}    {{.FieldType}}    `gorm:"column:{{.TableFieldName}}" json:"{{.TableFieldName}}" comment:"{{.TableFieldComment}}"`
    {{- end}}
}

func New() *{{.StructName}} {
	{{- if and .HasCreateTime .HasUpdateTime}}
	tn := time.Now()
    return &{{.StructName}}{CreateTime: tn, UpdateTime: tn}
    {{- else if .HasCreateTime}}
	t := time.Now()
    return &{{.StructName}}{CreateTime: tn}
    {{- else if .HasUpdateTime}}
	t := time.Now()
    return &{{.StructName}}{UpdateTime: tn}
	{{- else}}
	return &{{.StructName}}{}
    {{- end}}

}

func (this *{{.StructName}}) IsEdit() bool {
	return this.isEdit
}

func (this *{{.StructName}}) SetIsEdit(isEdit bool) {
	this.isEdit = isEdit
}

func (this *{{.StructName}}) SetFlag(flag higo.Flag) {
	this.currentFlag = flag
	this.isEdit = true
}

func (this *{{.StructName}}) Flag() higo.Flag {
	return this.currentFlag
}

func (this *{{.StructName}}) PriEmpty() bool {
	return this.{{.PrimaryId}} == 0
}
