package {{.Package}}

import (
	//"gitee.com/dengpju/higo-code/code"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-ioc/injector"
)

type {{.ModelImpl}} struct {
	*higo.Orm `inject:"Bean.NewOrm()"`
	{{range _, $v := .Fields}}
    {{end}}
	Id        int    `gorm:"column:id" json:"id"`
	Uname     string `gorm:"column:uname" json:"uname"`
	Utel      string `gorm:"column:u_tel" json:"utel"`
}

func init() {
	//init Validator
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

func (this *{{.ModelImpl}}) RegisterValidator() {
	//The custom tag, binding tag eg: binding:"custom_tag_name"
	//require import "gitee.com/dengpju/higo-code/code"
	// example
	//higo.RegisterValid(this).
	//	Tag("custom_tag_name",
	//		higo.Rule("required", code.Message("20000@custom_message")),
	//		higo.Rule("min=5", code.Message("20000@custom_message")))
}