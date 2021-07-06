package {{.Package}}

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/higo/request"
	"github.com/dengpju/higo-gin/higo/responser"
	"github.com/gin-gonic/gin"
)

type {{.Name}} struct {
    A string
}

func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{}
}

func (this *{{.Name}}) New() higo.IClass {
	return New{{.Name}}()
}

func (this *{{.Name}}) Route(hg *higo.Higo) {
    //TODO::example
	//route example
	hg.Get("/relative1", this.Example1, hg.Flag("this.Example1"), hg.Desc("Example1"))
	hg.Get("/relative2", this.Example2, hg.Flag("this.Example2"), hg.Desc("Example2"))
	hg.Get("/relative3", this.Example3, hg.Flag("this.Example3"), hg.Desc("Example3"))
	hg.Get("/relative4", this.Example4, hg.Flag("this.Example4"), hg.Desc("Example4"))
	hg.Get("/relative5", this.Example5, hg.Flag("this.Example5"), hg.Desc("Example5"))
    //route group example
    hg.AddGroup("/group_prefix", func() {
    	hg.Get("/relative6", this.Example6, hg.Flag("this.Example6"), hg.Desc("Example6"))
    })
}

func (this *{{.Name}}) Example1() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}

//responser string
func (this *{{.Name}}) Example2() string {
    //TODO::example code
    ctx := request.Context()
    //get parameter
    name := ctx.Query("name")
    return name
}

//responser interface{}
func (this *{{.Name}}) Example3() interface{} {
    //TODO::example code
    ctx := request.Context()
    //get parameter
    name := ctx.Query("name")
    return name
}

//example Model
type {{.Name}}ExampleModel struct {
	Id   int
	Name string
}

func (this *{{.Name}}ExampleModel) New() higo.IClass {
	return &{{.Name}}ExampleModel{}
}

func (this *{{.Name}}ExampleModel) Mutate(attrs ...higo.Property) higo.Model {
	higo.Propertys(attrs).Apply(this)
	return this
}

//responser Model
func (this *{{.Name}}) Example4(ctx *gin.Context) higo.Model {
    //TODO::example code
	model := &{{.Name}}ExampleModel{Id: 1, Name: "foo"}
	return model
}

//responser Models
func (this *{{.Name}}) Example5(ctx *gin.Context) higo.Models {
    //TODO::example code
	var models []*{{.Name}}ExampleModel
	models = append(models, &{{.Name}}ExampleModel{Id: 1, Name: "foo"}, &{{.Name}}ExampleModel{Id: 2, Name: "bar"})
	return higo.MakeModels(models)
}

//responser Json
func (this *{{.Name}}) Example6(ctx *gin.Context) higo.Json {
    //TODO::example code
	var models []*{{.Name}}ExampleModel
	models = append(models, &{{.Name}}ExampleModel{Id: 1, Name: "foo"}, &{{.Name}}ExampleModel{Id: 2, Name: "bar"})
	return models
}