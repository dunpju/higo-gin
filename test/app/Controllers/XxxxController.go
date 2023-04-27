package Controllers

import (
	"github.com/dengpju/higo-gin/higo"
	//"github.com/dengpju/higo-gin/higo/request"
	//"github.com/dengpju/higo-gin/higo/responser"
	//"github.com/gin-gonic/gin"
)

type XxxxController struct {
}

func NewXxxxController() *XxxxController {
	return &XxxxController{}
}

func (this *XxxxController) New() higo.IClass {
	return NewXxxxController()
}

func (this *XxxxController) Route(hg *higo.Higo) {
    /**
    //TODO::example
	//route example
	hg.Get("/relative1", this.Example1, hg.Flag("XxxxController.Example1"), hg.Desc("Example1"))
	hg.Get("/relative2", this.Example2, hg.Flag("XxxxController.Example2"), hg.Desc("Example2"))
	hg.Get("/relative3", this.Example3, hg.Flag("XxxxController.Example3"), hg.Desc("Example3"))
	hg.Get("/relative4", this.Example4, hg.Flag("XxxxController.Example4"), hg.Desc("Example4"))
	hg.Get("/relative5", this.Example5, hg.Flag("XxxxController.Example5"), hg.Desc("Example5"))
    //route group example
    hg.AddGroup("/group_prefix", func() {
        hg.Get("/relative6", this.Example6, hg.Flag("XxxxController.Example6"), hg.Desc("Example6"))
    	hg.Get("/list", this.List, hg.Flag("XxxxController.List"), hg.Desc("List"))
    	hg.Post("/add", this.Add, hg.Flag("XxxxController.Add"), hg.Desc("Add"))
    	hg.Put("/edit", this.Edit, hg.Flag("XxxxController.Edit"), hg.Desc("Edit"))
    	hg.Delete("/delete", this.Delete, hg.Flag("XxxxController.Delete"), hg.Desc("Delete"))
    })
    */
}

/**
func (this *XxxxController) List() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *XxxxController) Add() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *XxxxController) Edit() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *XxxxController) Delete() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *XxxxController) Example1() {
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
func (this *XxxxController) Example2() string {
    //TODO::example code
    ctx := request.Context()
    //get parameter
    name := ctx.Query("name")
    return name
}
*/

/**
//responser interface{}
func (this *XxxxController) Example3() interface{} {
    //TODO::example code
    ctx := request.Context()
    //get parameter
    name := ctx.Query("name")
    return name
}
*/

/**
//example Model
type XxxxControllerModel struct {
	Id   int
	Name string
}

func (this *XxxxControllerModel) New() higo.IClass {
	return &XxxxControllerModel{}
}

func (this *XxxxControllerModel) Mutate(attrs ...higo.Property) higo.Model {
	higo.Propertys(attrs).Apply(this)
	return this
}
*/

/**
//responser Model
func (this *XxxxController) Example4(ctx *gin.Context) higo.Model {
    //TODO::example code
	model := &XxxxControllerModel{Id: 1, Name: "foo"}
	return model
}
*/

/**
//responser Models
func (this *XxxxController) Example5(ctx *gin.Context) higo.Models {
    //TODO::example code
	var models []*XxxxControllerModel
	models = append(models, &XxxxControllerModel{Id: 1, Name: "foo"}, &XxxxControllerModel{Id: 2, Name: "bar"})
	return higo.MakeModels(models)
}
*/

/**
//responser Json
func (this *XxxxController) Example6(ctx *gin.Context) higo.Json {
    //TODO::example code
	var models []*XxxxControllerModel
	models = append(models, &XxxxControllerModel{Id: 1, Name: "foo"}, &XxxxControllerModel{Id: 2, Name: "bar"})
	return models
}
*/