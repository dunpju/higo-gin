package {{.Package}}

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

type {{.ModelImpl}} struct {
	*higo.Orm `inject:"Bean.NewOrm()"`
	{{- range .TplFields}}
	{{.Field}}        {{.Type}}    `gorm:"column:{{.DbField}}" json:"{{.DbField}}" comment:"{{.Comment}}"`
    {{- end}}
}

var (
    Models []*{{.ModelImpl}}
)

//init Validator
func init() {
	New().RegisterValidator()
}

func New(attrs ...higo.Property) *{{.ModelImpl}} {
	impl := &{{.ModelImpl}}{}
	higo.Propertys(attrs).Apply(impl)
	injector.BeanFactory.Apply(impl)
	return impl
}

func (this *{{.ModelImpl}}) New() higo.IClass {
	return New()
}

func (this *{{.ModelImpl}}) TableName() string {
	return "{{.TableName}}"
}

func (this *{{.ModelImpl}}) Mutate(attrs ...higo.Property) higo.Model {
	higo.Propertys(attrs).Apply(this)
	return this
}

//The custom tag, binding the tag eg: binding:"custom_tag_name"
//require import "gitee.com/dengpju/higo-code/code"
//
//example code:
//func (this *ModelImpl) RegisterValidator() {
//	higo.RegisterValid(this).
//		Tag("custom_tag_name",
//			higo.Rule("required", code.Message("20000@custom_message")),
//			higo.Rule("min=5", code.Message("20000@custom_message")))
//}
func (this *ModelImpl) RegisterValidator() {
}

func (this *ModelImpl) GetByID(ID {{.PriType}}, columns ...string) {
	this.Mapper(squirrel.Select(columns...).From(this.TableName()).Where("`{{.PRI}}` = ?", ID).ToSql()).Query().Scan(&this)
}

func (this *ModelImpl) GetByIDS(IDS []string, columns ...string) *gorm.DB {
	return this.Mapper(squirrel.Select(columns...).From(this.TableName()).Where("`{{.PRI}}` IN(?)", strings.Join(IDS, ",")).ToSql()).Query()
}