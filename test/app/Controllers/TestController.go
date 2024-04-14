package Controllers

import (
	"github.com/dunpju/higo-gin/higo"
	//"github.com/dunpju/higo-gin/higo/request"
	//"github.com/dunpju/higo-gin/higo/responser"
	//"github.com/gin-gonic/gin"
)

const SelfTestController = "app/Controllers/TestController"

type TestController struct {
}

func NewTestController() *TestController {
	return &TestController{}
}

func (this *TestController) New() higo.IClass {
	return NewTestController()
}

func (this *TestController) Route(hg *higo.Higo) {
    /**
    //TODO::example
	//route example
	hg.Get("/relative1", this.Example1, hg.Flag("TestController.Example1"), hg.Desc("Example1"))
	hg.Get("/relative2", this.Example2, hg.Flag("TestController.Example2"), hg.Desc("Example2"))
	hg.Get("/relative3", this.Example3, hg.Flag("TestController.Example3"), hg.Desc("Example3"))
	hg.Get("/relative4", this.Example4, hg.Flag("TestController.Example4"), hg.Desc("Example4"))
	hg.Get("/relative5", this.Example5, hg.Flag("TestController.Example5"), hg.Desc("Example5"))
    //route group example
    hg.AddGroup("/group_prefix", func() {
        hg.Get("/relative6", this.Example6, hg.Flag("TestController.Example6"), hg.Desc("Example6"))
    	hg.Get("/list", this.List, hg.Flag("TestController.List"), hg.Desc("List"))
    	hg.Post("/add", this.Add, hg.Flag("TestController.Add"), hg.Desc("Add"))
    	hg.Put("/edit", this.Edit, hg.Flag("TestController.Edit"), hg.Desc("Edit"))
    	hg.Delete("/delete", this.Delete, hg.Flag("TestController.Delete"), hg.Desc("Delete"))
    })
    */
}

/**
func (this *TestController) List() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *TestController) Add() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *TestController) Edit() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *TestController) Delete() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *TestController) Example1() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
//responser string
func (this *TestController) Example2() string {
    //TODO::example code
    ctx := request.Context()
    //get parameter
    name := ctx.Query("name")
    return name
}
*/

/**
//responser interface{}
func (this *TestController) Example3() interface{} {
    //TODO::example code
    ctx := request.Context()
    //get parameter
    name := ctx.Query("name")
    return name
}
*/

/**
//example Model
type TestControllerModel struct {
	Id   int
	Name string
}

func (this *TestControllerModel) New() higo.IClass {
	return &TestControllerModel{}
}

func (this *TestControllerModel) Mutate(attrs ...higo.Property) higo.Model {
	higo.Propertys(attrs).Apply(this)
	return this
}
*/

/**
//responser Model
func (this *TestController) Example4(ctx *gin.Context) higo.Model {
    //TODO::example code
	model := &TestControllerModel{Id: 1, Name: "foo"}
	return model
}
*/

/**
//responser Models
func (this *TestController) Example5(ctx *gin.Context) higo.Models {
    //TODO::example code
	var models []*TestControllerModel
	models = append(models, &TestControllerModel{Id: 1, Name: "foo"}, &TestControllerModel{Id: 2, Name: "bar"})
	return higo.MakeModels(models)
}
*/

/**
//responser Json
func (this *TestController) Example6(ctx *gin.Context) higo.Json {
    //TODO::example code
	var models []*TestControllerModel
	models = append(models, &TestControllerModel{Id: 1, Name: "foo"}, &TestControllerModel{Id: 2, Name: "bar"})
	return models
}
*/