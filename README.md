# higo-gin
基于gin封装的脚手架

## <span id="top">功能说明</span>
[控制器](#controller)、[简易依赖注入](#inject)、[中间件](#middleware)、[表达式](#express)、[任务组件](#task)、[限流器](#limiter)、[开发者工具](#devtool)等

## 目录结构
仅参考:
```
project
  |--app
  |   |--beans
  |   |--controllers
  |   |--dao
  |   |--entity
  |   |--enums
  |   |--errcode
  |   |--exception
  |   |--middlewares
  |   |--models
  |   |--params
  |   |--services
  |   |--utils
  |--bin
  |   |--yaml
  |   |   |--20000.yaml
  |   |--enum_cmd.md
  |   |--main.go
  |--env
  |   |--app.yaml
  |   |--database.yaml
  |   |--serve.yaml
  |--router
  |   |--http.go
  |--go.mod
```

## 安装
```
go get -u github.com/dunpju/higo-gin@v1.0.0
```
# 使用
### 配置
配置都放在env目录下,在启动服务时会自动加载yaml配置
###### app.yaml文件
```
# app config目录
APP_CONFIG: /env/config
```
###### database.yaml文件
```
DB:
  Default:
    Driver: "mysql"
    Host: "127.0.0.1"
    Port: "3306"
    Database: "db_name"
    Username: "root"
    Password: ""
    Charset: "utf8mb4"
    Collation: "utf8mb4_unicode_ci"
    Prefix: ""
    LogMode: true
    MaxIdle: 200
    MaxOpen: 1
    MaxLifetime: 10

Redis:
  Default:
    Host: ""
    Auth: ""
    Port: 6379
    Db: 0
    Pool:
      Min_Connections: 1
      Max_Connections: 10
      Connect_Timeout: 10.0
      Wait_Timeout: 3.0
      Heartbeat: -1
      Max_Idle: 3
      Max_Idle_Time: 60
      Max_Conn_Lifetime: 10
      Wait: true
```
###### serve.yaml文件
```
# Http服务器
HTTP_HOST:
  Type: http
  Name: http
  Addr: 0.0.0.0:1254
  ReadTimeout: 5
  WriteTimeout: 10

# Https服务器
HTTPS_HOST:
  Type: http
  Name: https
  Addr: 0.0.0.0:1255
  ReadTimeout: 5
  WriteTimeout: 10

# Websocket服务器
WEBSOCKET_HOST:
  Type: websocket
  Name: websocket
  Addr: 0.0.0.0:6125
  ReadTimeout: 5
  WriteTimeout: 10
```
###### 获取配置
```
获取all: config.All()
获取key: config.Get("key") 或者 config.Get("env.key")
获取env: config.Env("key")
获取app.yaml: config.App("key")
获取anno.yaml: config.Anno("key")
```

### 路由
路由可以统一写在router/http.go文件里
```
// https api 接口
type Https struct {
    *higo.Serve `inject:"Bean.NewServe('env.serve.HTTPS_HOST')"`
}

func NewHttps() *Https {
    return &Https{}
}

func (this *Https) Static(ctx *gin.Context) bool {
    ok1, err := regexp.MatchString("^/index/", ctx.Request.URL.Path)
    if nil != err {
        panic(err)
    }
    ok2, err := regexp.MatchString(`.(js|css|woff|ttf|ico|png)$`, ctx.Request.URL.Path)
    if nil != err {
        panic(err)
    }
    if ok1 || ok2 {
        return true
    }
    return false
}

// 路由装载器
func (this *Https) Loader(hg *higo.Higo) {
    // 静态文件
    hg.Static("/index/", "dist")
    hg.Static("/static/", "dist/static")
    hg.StaticFile("/favicon.ico", "dist/favicon.ico")
    this.Api(hg)
}

// api 路由
func (this *Https) Api(hg *higo.Higo) {
    hg.AddGroup("/v3", func() {
        hg.Get("/test1", HttpsTestThrow2)
    })
}
```
或者写在controller类Route方法里
```
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
```
### 启动
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
### 鉴权
鉴权需要中间件和路由配合,默认情况所有Api都是需要鉴权的,在注册路由时可以router.IsAuth(false)取消鉴权```hg.Post("/test1", this.HttpsTestValidate, router.IsAuth(false))```
###### 中间件Auth.go
```
// 鉴权
type Auth struct{}

// 构造函数
func NewAuth() *Auth {
    return &Auth{}
}

func (this *Auth) Middle(hg *higo.Higo) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if !router.NewHttps().Static(ctx) {
            if route, ok := hg.GetRoute(ctx.Request.Method, ctx.Request.URL.Path); ok {
                if route.IsAuth() && !route.IsStatic() {
                    token := ctx.GetHeader("Authorization")
                    if "" == token {
                        exception.Throw(exception.Message(errcode.TokenEmpty.Message()), exception.Code(int(errcode.TokenEmpty)))
                    }
                    // token解析
                }
            } else {
                exception.Throw(exception.Message(errcode.NotFound.Message()), exception.Code(int(errcode.NotFound)))
            }
        }
    }
}
```

### <span id="controller">控制器</span> <font size=1>[top](#top)</font>
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
    "github.com/dunpju/higo-gin/higo"
    "github.com/dunpju/higo-gin/higo/request"
    "github.com/dunpju/higo-gin/higo/responser"
    "github.com/dunpju/higo-router/router"
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
	ctx := request.Context()
    name := ctx.Query("name")
    responser.SuccessJson("success", 10000, name)
}

func (this *AdminController) Add() {
    ctx := request.Context()
    name := ctx.Query("name")
    responser.SuccessJson("success", 10000, name)
}

func (this *AdminController) Edit() {
    ctx := request.Context()
    name := ctx.Query("name")
    responser.SuccessJson("success", 10000, name)
}

func (this *AdminController) Delete() {
    ctx := request.Context()
    name := ctx.Query("name")
    responser.SuccessJson("success", 10000, name)
}

func (this *AdminController) Example1() {
    ctx := request.Context()
    name := ctx.Query("name")
    responser.SuccessJson("success", 10000, name)
}

//responser string
func (this *AdminController) Example2() string {
    ctx := request.Context()
    name := ctx.Query("name")
    return name
}

//responser interface{}
func (this *AdminController) Example3() interface{} {
    ctx := request.Context()
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
    model := &AdminControllerModel{Id: 1, Name: "foo"}
    return model
}

//responser Models
func (this *AdminController) Example5(ctx *gin.Context) higo.Models {
    var models []*AdminControllerModel
    models = append(models, &AdminControllerModel{Id: 1, Name: "foo"}, &AdminControllerModel{Id: 2, Name: "bar"})
    return higo.MakeModels(models)
}

//responser Json
func (this *AdminController) Example6(ctx *gin.Context) higo.Json {
    var models []*AdminControllerModel
    models = append(models, &AdminControllerModel{Id: 1, Name: "foo"}, &AdminControllerModel{Id: 2, Name: "bar"})
    return models
}
```

###### <span id="Websocket">Websocket</span> <font size=1>[top](#top)</font>
```
type WebsocketController struct {
}

func NewWebsocketController() *WebsocketController {
	return &WebsocketController{}
}

func (this *WebsocketController) New() higo.IClass {
	return NewWebsocketController()
}

func (this *WebsocketController) Route(hg *higo.Higo) {
	hg.Ws("/conn", this.Conn, hg.Desc("conn"))
	hg.Ws("/echo", this.Echo, hg.Flag("WebsocketController.Echo"), hg.Desc("Echo"))
	hg.Ws("/send_all", this.SendAll, hg.Flag("WebsocketController.SendAll"), hg.Desc("SendAll"))
}

// webSocket请求
func (this *WebsocketController) Conn(ctx *gin.Context) higo.WsWriteMessage {
	loginEntity := Entity.NewLoginEntity()
	err := ctx.ShouldBind(loginEntity)
	if err != nil {
		panic(err)
	}
	fmt.Println("Conn", loginEntity)

	return higo.WsRespStruct(loginEntity)
}

func (this *WebsocketController) Echo(ctx *gin.Context) higo.WsWriteMessage {
	return higo.WsRespString("echo")
}

func (this *WebsocketController) SendAll(ctx *gin.Context) string {
	higo.WsContainer.SendAll(ctx.Query("msg"))
	return "ok"
}
```

###### <span id="redis">Redis</span> <font size=1>[top](#top)</font>
redis工具基于github.com/gomodule/redigo/redis封装<br>
启动时开启redis连接池
```
func main() {
    higo.Init(sliceutil.NewSliceString(".", "test", "")).
        IsRedisPool(). // 开启redis连接池
        Boot()
}
```
使用
```
import (
    "fmt"
    "github.com/dunpju/higo-gin/higo"
)

func Test() {
    higo.Redis.Set("name", rand.Intn(1000))
    v := higo.Redis.Get("name")
    fmt.Println(v)
    // 获取github.com/gomodule/redigo/redis连接
    higo.Redis.Conn()            
}
```

### <span id="inject">依赖注入</span> <font size=1>[top](#top)</font>
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

### <span id="middleware">中间件</span> <font size=1>[top](#top)</font>
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

### <span id="express">表达式</span> <font size=1>[top](#top)</font>
表达式封装了独立的包```github.com/dunpju/higo-express/express```
```
// 执行表达式
express.Run("MyBean.DemoService()")
```

### <span id="task">任务组件</span> <font size=1>[top](#top)</font>
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

### <span id="limiter">限流器</span> <font size=1>[top](#top)</font>
使用令牌桶算法生成token，并实现LRU算法的淘汰机制;
###### 使用限流器
```
hg.Get("/test1", ratelimit.Limiter(3, 1)(Test))
说明：ratelimit.Limiter(3, 1)表示令牌桶容量为3,每秒产生1个令牌
```

### <span id="devtool">开发者工具</span> <font size=1>[top](#top)</font>
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
