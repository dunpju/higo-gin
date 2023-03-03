# higo-gin
基于gin封装的脚手架

## 安装
go get -u github.com/dengpju/higo-gin@v1.1.54

## 启动
```
func main() {
	beanConfig := Beans.NewMyBean()

	higo.Init(sliceutil.NewSliceString(".", "test", "")). // 指定root目录
		Middleware(Middlewares.NewCors(), Middlewares.NewRunLog()). // 注册中间件
		AddServe(router.NewHttp(), Middlewares.NewHttp()). // 添加http服务、服务中间件
		AddServe(router.NewHttps(), beanConfig). // 添加https服务、服务bean
		AddServe(router.NewWebsocket()).// 添加websocket服务
		Beans(beanConfig). // 全局bean
		Cron("0/3 * * * * *", func() { // 添加定时任务
		    log.Println("3秒执行一次")
		}).
		Boot()
}
```

## 功能说明
控制器、简易依赖注入、中间件、表达式、任务组件、开发者工具等

### 控制器
控制器代码可以通过开发者工具指令直接生成基础代码,并且采用AST自动注册到Bean实例内。

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
注入对象,需要先注册到Bean实例内

###### example
```
type MyBean struct{ higo.Bean }

func NewMyBean() *MyBean {
	return &MyBean{}
}

// 手动添加配置进行注册
func (this *MyBean) DemoService() *Services.DemoService {
	return Services.NewDemoService()
}
```
###### 使用
```
type DemoController struct {
    // multiple
	DemoService1 *Services.DemoService `inject:"MyBean.DemoService()"`
	// single
	DemoService2 *Services.DemoService `inject:"-"`
	Name string
}
```
###### Value注入
Value注入需要在```env/Config/anno.yaml```文件中添加配置
```
"value":
  user:
    score: "100"
    age: "19"
```
###### 使用
```
type DemoController struct {
    Age  *anno.Value  `prefix:"user.age"`
}
```

### 中间件
只要实现IMiddleware接口都被认为是中间件
###### example
```
type Cors struct {}

func NewCors() *Cors {
	return &Cors{}
}

func (this *Cors) Middle(hg *higo.Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		method := cxt.Request.Method
		origin := cxt.Request.Header.Get("Origin")
		if origin != "" {
			cxt.Header("Access-Control-Allow-Origin", "*")
			cxt.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			cxt.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			cxt.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			cxt.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			cxt.AbortWithStatus(http.StatusNoContent)
		}

		cxt.Next()
	}
}
```

### 表达式
表达式封装了独立的包```github.com/dengpju/higo-express/express```
```
// 执行表达式
express.Run("MyBean.DemoService()")
```

### 任务组件
使用定时任务
```
func main() {
	higo.Init(sliceutil.NewSliceString(".", "test", "")).
		Cron("0/3 * * * * *", func() {
		    log.Println("3秒执行一次")
		}).
		Boot()
}
```

### 开发者工具
###### 构建Controller
```
go run bin\main.go -gen=controller -out=app\controllers -name=controller_name
```
###### 构建Model
```
go run bin\main.go -gen=model -out=app\models -name=table_name
table_name如果是all,就构建所有表;
构建model时会自动构建Dao层、Entity层
```
###### 构建Enum
```
go run bin\main.go -gen=enum -out=app\enums -name="-e=flag_state -f=状态:enable-1-启用,disabled-2-禁用"
或者将指令写入文件,例如将上面指令写入enum_cmd.md文件,然后执行
go run bin\main.go -gen=enum -out=app\enums -name=bin\enum_cmd.md
```
###### 构建Service
```
go run bin\main.go -gen=service -out=app\service -name=TestService
```
###### 构建Code
bin\yaml\20000.yaml配置
```
success:
  code: "20000"
  message: "成功"
  iota: "yes"
nonexistence:
```
执行下面指令:
```
go run bin\main.go -gen=code -name=ErrorCode -out=app\errcode -path=bin\yaml -auto=yes -force=yes
```
###### 构建Param
```
go run bin\main.go -gen=param -out=app\params -name=DemoList
```
自定义参数校验:
```
type DemoList struct {
    // binding:"table",table是自定义tag
	Table uint64 `form:"table" binding:"table"` // get from the form
}

func NewParamDemoList(ctx *gin.Context) *DemoList {
	param := &DemoList{}
	higo.Validate(param).Receiver(ctx.ShouldBindQuery(param)).Unwrap() // get from the form
	return param
}

// RegisterValidator The custom tag, binding the tag eg: binding:"custom_tag_name"
// require import "gitee.com/dengpju/higo-code/code"
//
//example code:
//func (this *DemoList) RegisterValidator() *higo.Verify {
//	return higo.RegisterValidator(this).
//		Tag("custom_tag_name",
//			higo.Rule("required", Codes.Success),
//			higo.Rule("min=5", Codes.Success))
//		Tag("custom_tag_name2",
//			higo.Rule("required", func() higo.ValidatorToFunc {
//              return func(fl validator.FieldLevel) (bool, code.ICode) {
//                  fmt.Println(fl.Field().Interface())
//                  return true, MinError
//              }
//          }()))
//  Or
//  return higo.Verifier() // Manual call Register Validate: higo.Validate(verifier)
//}
func (this *DemoList) RegisterValidator() *higo.Verify {
	return higo.Verifier().
	    // 自定义tag
		Tag("table", higo.Rule("required", func() higo.ValidatorToFunc {
			return func(fieldLevel validator.FieldLevel) (bool, code.ICode) {
			    // 手动校验
				err := EnumTable.Inspect(int(fieldLevel.Field().Interface().(uint64)))
				if err != nil {
				    // 未通过
					return false, errcode.Chain
				}
				// 通过
				return true, nil
			}
		}()))
}
```
###### 组合校验规则
```
type Add struct {
	Name      string  `json:"name" binding:"name"`
	State     int     `json:"state" binding:"state"`
	Remark    string  `json:"remark"`
}

type Delete struct {
	Id int64 `json:"id" binding:"id"`
}

type Edit struct {
	Add
	Delete
}

func NewEdit(ctx *gin.Context) *Edit {
	param := &Edit{}
	higo.Validate(param).Receiver(ctx.ShouldBindBodyWith(param, binding.JSON)).Unwrap() // get from the json multiple binding
	return param
}

func (this *Edit) RegisterValidator() *higo.Verify {
	return higo.Verifier().Use(&Add{}, &Delete{})
}
```
使用AST技术实现Code、Controller、Dao、Entity、Model、Enum、Service、Validate等生成器，
定时器、IOC容器、Bean工厂、限流器(采用LRU算法)、路由收集器(采用前缀树)、Websocket调度器、
自动校验器、Code容器、Enum容器、基于gorm封装的Orm工具、事件处理机制、协程字节流处理工具、
动态参数工具(模拟多态)、使用Protobuf实现简单的存储配合coroutine实现快速的CRUD功能;