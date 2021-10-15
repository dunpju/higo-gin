package {{.PackageName}}

import (
	//"gitee.com/dengpju/higo-code/code"
	"github.com/Masterminds/squirrel"
    "github.com/dengpju/higo-gin/higo"
    "github.com/dengpju/higo-ioc/injector"
    "github.com/jinzhu/gorm"
    "strings"
	{{- range $impo := .Imports}}
    {{$impo}}
    {{- end}}
)

type {{.StructName}} struct {
	*higo.Orm    `inject:"Bean.NewOrm()"`
	{{- range .StructFields}}
	{{.FieldName}}    {{.FieldType}}    `gorm:"column:{{.TableFieldName}}" json:"{{.TableFieldName}}" comment:"{{.TableFieldComment}}"`
    {{- end}}
}

//init Validator
func init() {
	New().RegisterValidator()
}

func New(attrs ...higo.Property) *{{.StructName}} {
	impl := &{{.StructName}}{}
	higo.Propertys(attrs).Apply(impl)
	injector.BeanFactory.Apply(impl)
	return impl
}

func (this *{{.StructName}}) New() higo.IClass {
	return New()
}

func (this *{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}

func (this *{{.StructName}}) Mutate(attrs ...higo.Property) higo.Model {
	higo.Propertys(attrs).Apply(this)
	return this
}

//The custom tag, binding the tag eg: binding:"custom_tag_name"
//require import "gitee.com/dengpju/higo-code/code"
//
//example code:
//func (this *StructName) RegisterValidator() *higo.Valid {
//	return higo.RegisterValid(this).
//		Tag("custom_tag_name",
//			higo.Rule("required", Codes.Success),
//			higo.Rule("min=5", Codes.Success))
//  Or
//  return higo.Verifier() // Manual call Register Validate: higo.Validate(verifier)
//}
func (this *{{.StructName}}) RegisterValidator() *higo.Valid {
    return higo.RegisterValid(this)
}

func (this *{{.StructName}}) Exist() bool {
	return this.{{.PrimaryId}} > 0
}

func (this *{{.StructName}}) GetBy{{.PrimaryId}}({{.SmallHumpPrimaryId}} {{.PrimaryIdType}}, columns ...string) *gorm.DB {
	return this.Mapper(squirrel.Select(columns...).From(this.TableName()).Where({{.PrimaryId}} + " = ?", {{.SmallHumpPrimaryId}}).ToSql()).Query()
}

func (this *{{.StructName}}) GetBy{{.PrimaryId}}s({{.SmallHumpPrimaryId}}s []string, columns ...string) *gorm.DB {
	return this.Mapper(squirrel.Select(columns...).From(this.TableName()).Where({{.PrimaryId}} + " IN(?)", strings.Join({{.SmallHumpPrimaryId}}s, ",")).ToSql()).Query()
}

func (this *{{.StructName}}) Paginate(perPage, page uint64) *higo.Pager {
	models := make([]*{{.StructName}}, 0)
	pager := higo.NewPager(perPage, page)
	this.Table(this.TableName()). /**Where().*/ Paginate(pager).Find(&models)
	pager.Items = models
	return pager
}