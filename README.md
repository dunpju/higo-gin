# higo-gin
基于gin封装的脚手架

## 安装
go get -u github.com/dengpju/higo-gin@v1.1.54

## 功能说明
控制器、简易依赖注入、中间件、表达式、任务组件、开发者工具等

### 控制器
控制器代码可以通过开发者工具指令直接生成基础代码,并且采用AST技术自动注册到BeanFactory内。

###### 构建指令:
```
go run bin\main.go -gen=controller -out=app\controllers -name=Admin
-out: 输出目录
-name: 控制器名称
```

###### example
```
package controllers

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/higo/request"
	"github.com/dengpju/higo-gin/higo/responser"
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-gin/higo/request"
	"github.com/dengpju/higo-gin/higo/responser"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
}

func NewAdminController() *AdminController {
	return &AdminController{}
}

func (this *AdminController) New() higo.IClass {
	return NewAdminController()
}

func (this *AdminController) Route(hg *higo.Higo) {
	//route example
	hg.Get("/relative1", this.Example1, hg.Flag("AdminController.Example1"), hg.Desc("Example1"))
	hg.Get("/relative2", this.Example2, hg.Flag("AdminController.Example2"), hg.Desc("Example2"))
	hg.Get("/relative3", this.Example3, hg.Flag("AdminController.Example3"), hg.Desc("Example3"))
	hg.Get("/relative4", this.Example4, hg.Flag("AdminController.Example4"), hg.Desc("Example4"))
	hg.Get("/relative5", this.Example5, hg.Flag("AdminController.Example5"), hg.Desc("Example5"))
	//route group example
	hg.AddGroup("/group_prefix", func() {
	    hg.Get("/relative6", this.Example6, hg.Flag("AdminController.Example6"), hg.Desc("Example6"))
	    hg.Get("/list", this.List, hg.Flag("AdminController.List"), hg.Desc("List"))
	    hg.Post("/add", this.Add, hg.Flag("AdminController.Add"), hg.Desc("Add"))
	    hg.Put("/edit", this.Edit, hg.Flag("AdminController.Edit"), hg.Desc("Edit"))
	    hg.Delete("/delete", this.Delete, hg.Flag("AdminController.Delete"), hg.Desc("Delete"))
	})
}

func (this *AdminController) List() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}

func (this *AdminController) Add() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}

func (this *AdminController) Edit() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}

func (this *AdminController) Delete() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}

func (this *AdminController) Example1() {
    //TODO::example code
	ctx := request.Context()
	//get parameter
    name := ctx.Query("name")
    //responser
    responser.SuccessJson("success", 10000, name)
}

//responser string
func (this *AdminController) Example2() string {
    //TODO::example code
    ctx := request.Context()
    //get parameter
    name := ctx.Query("name")
    return name
}

//responser interface{}
func (this *AdminController) Example3() interface{} {
    //TODO::example code
    ctx := request.Context()
    //get parameter
    name := ctx.Query("name")
    return name
}

//example Model
type AdminControllerModel struct {
	Id   int
	Name string
}

func (this *AdminControllerModel) New() higo.IClass {
	return &AdminControllerModel{}
}

func (this *AdminControllerModel) Mutate(attrs ...higo.Property) higo.Model {
	higo.Propertys(attrs).Apply(this)
	return this
}

//responser Model
func (this *AdminController) Example4(ctx *gin.Context) higo.Model {
    //TODO::example code
	model := &AdminControllerModel{Id: 1, Name: "foo"}
	return model
}

//responser Models
func (this *AdminController) Example5(ctx *gin.Context) higo.Models {
    //TODO::example code
	var models []*AdminControllerModel
	models = append(models, &AdminControllerModel{Id: 1, Name: "foo"}, &AdminControllerModel{Id: 2, Name: "bar"})
	return higo.MakeModels(models)
}

//responser Json
func (this *AdminController) Example6(ctx *gin.Context) higo.Json {
    //TODO::example code
	var models []*AdminControllerModel
	models = append(models, &AdminControllerModel{Id: 1, Name: "foo"}, &AdminControllerModel{Id: 2, Name: "bar"})
	return models
}

```


### 依赖注入

### 中间件

### 表达式

### 任务组件

### 开发者工具
使用AST技术实现Code、Controller、Dao、Entity、Model、Enum、Service、Validate等生成器，
定时器、IOC容器、Bean工厂、限流器(采用LRU算法)、路由收集器(采用前缀树)、Websocket调度器、自动校验器、Code容器、Enum容器、基于gorm封装的Orm工具、事件处理机制、协程字节流处理工具、动态参数工具(模拟多态)、使用Protobuf实现简单的存储配合coroutine实现快速的CRUD功能;