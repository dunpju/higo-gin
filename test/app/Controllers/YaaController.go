package Controllers

import (
	"github.com/dengpju/higo-gin/higo"
	//"github.com/dengpju/higo-gin/higo/request"
	//"github.com/dengpju/higo-gin/higo/responser"
	//"github.com/gin-gonic/gin"
)

type YaaController struct {
}

func NewYaaController() *YaaController {
	return &YaaController{}
}

func (this *YaaController) New() higo.IClass {
	return NewYaaController()
}

func (this *YaaController) Route(hg *higo.Higo) {
    /**
    //TODO::example
	//route example
	hg.Get("/relative1", this.Example1, hg.Flag("YaaController.Example1"), hg.Desc("Example1"))
	hg.Get("/relative2", this.Example2, hg.Flag("YaaController.Example2"), hg.Desc("Example2"))
	hg.Get("/relative3", this.Example3, hg.Flag("YaaController.Example3"), hg.Desc("Example3"))
	hg.Get("/relative4", this.Example4, hg.Flag("YaaController.Example4"), hg.Desc("Example4"))
	hg.Get("/relative5", this.Example5, hg.Flag("YaaController.Example5"), hg.Desc("Example5"))
    //route group example
    hg.AddGroup("/group_prefix", func() {
        hg.Get("/relative6", this.Example6, hg.Flag("YaaController.Example6"), hg.Desc("Example6"))
    	hg.Get("/list", this.List, hg.Flag("YaaController.List"), hg.Desc("List"))
    	hg.Post("/add", this.Add, hg.Flag("YaaController.Add"), hg.Desc("Add"))
    	hg.Put("/edit", this.Edit, hg.Flag("YaaController.Edit"), hg.Desc("Edit"))
    	hg.Delete("/delete", this.Delete, hg.Flag("YaaController.Delete"), hg.Desc("Delete"))
    })
    */
}

/**
func (this *YaaController) List() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *YaaController) Add() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *YaaController) Edit() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *YaaController) Delete() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}
*/

/**
func (this *YaaController) Example1() {
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
func (this *YaaController) Example2() string {
    //TODO::example code
    ctx := request.Context()
    //get parameter
    name := ctx.Query("name")
    return name
}
*/

/**
//responser interface{}
func (this *YaaController) Example3() interface{} {
    //TODO::example code
    ctx := request.Context()
    //get parameter
    name := ctx.Query("name")
    return name
}
*/

/**
//example Model
type YaaControllerModel struct {
	Id   int
	Name string
}

func (this *YaaControllerModel) New() higo.IClass {
	return &YaaControllerModel{}
}

func (this *YaaControllerModel) Mutate(attrs ...higo.Property) higo.Model {
	higo.Propertys(attrs).Apply(this)
	return this
}
*/

/**
//responser Model
func (this *YaaController) Example4(ctx *gin.Context) higo.Model {
    //TODO::example code
	model := &YaaControllerModel{Id: 1, Name: "foo"}
	return model
}
*/

/**
//responser Models
func (this *YaaController) Example5(ctx *gin.Context) higo.Models {
    //TODO::example code
	var models []*YaaControllerModel
	models = append(models, &YaaControllerModel{Id: 1, Name: "foo"}, &YaaControllerModel{Id: 2, Name: "bar"})
	return higo.MakeModels(models)
}
*/

/**
//responser Json
func (this *YaaController) Example6(ctx *gin.Context) higo.Json {
    //TODO::example code
	var models []*YaaControllerModel
	models = append(models, &YaaControllerModel{Id: 1, Name: "foo"}, &YaaControllerModel{Id: 2, Name: "bar"})
	return models
}
*/