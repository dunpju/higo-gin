package V3

import (
	"fmt"
	"gitee.com/dengpju/higo-code/code"
	"github.com/dengpju/higo-annotation/anno"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/higo/request"
	"github.com/dengpju/higo-gin/test/app/Controllers"
	"github.com/dengpju/higo-gin/test/app/Exception"
	"github.com/dengpju/higo-gin/test/app/Models/UserModel"
	"github.com/dengpju/higo-gin/test/app/Services"
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-throw/exception"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"time"
)

type DemoController struct {
	Age         *anno.Value           `prefix:"user.age"`
	DemoService *Services.DemoService `inject:"MyBean.DemoService()"`
	//*higo.Orm   `inject:"Bean.NewOrm()"`
	//*redis.Pool `inject:"Bean.NewRedisPool()"`
	Name string
}

func NewDemoController() *DemoController {
	return &DemoController{}
}

func (this *DemoController) New() higo.IClass {
	return NewDemoController()
}

func (this *DemoController) Route(hg *higo.Higo) {
	hg.AddGroup("/https/v3", func() {
		//hg.AddGroup("/user", func() {
		//	hg.Post("/login", this.Login, hg.Flag("Login"), hg.Desc("V3 登录"))
		//})
		hg.Get("/test_throw2132", HttpsTestThrow2, hg.Flag("TestThrow"), hg.Desc("V3 测试异常111"))
		hg.Get("/test_throw21321", this.HttpsTestThrow2, hg.Flag("TestThrow"), hg.Desc("V3 测试异常111"))
		hg.Post("/test_get21111132", this.HttpsTestGet, hg.Flag("TestGet"), hg.Desc("V3 测试GET"))
		hg.Post("/test_validator2132", this.HttpsTestValidate, hg.Flag("HttpsTestValidate"), hg.Desc("V3 测试校验器"), router.IsAuth(false))
	})
}

func HttpsTestThrow2(ctx *gin.Context)  {
	fmt.Println("ggg")
	fmt.Printf("%p\n", higo.Di("test/app/Controllers/WebsocketController"))
	w1 := higo.Di("test/app/Controllers/WebsocketController")
	w1.(*Controllers.WebsocketController).K = "tt"
	w2 := higo.Di("test/app/Controllers/WebsocketController")
	w2.(*Controllers.WebsocketController).K = "tt1"
	fmt.Println(w1, w2)
}

// 测试异常
func (this *DemoController) HttpsTestThrow2(ctx *gin.Context) string {
	fmt.Println(ctx.Query("id"))
	fmt.Println(111)
	fmt.Println(&this)
	fmt.Println(this.Age)
	var s []map[string]interface{}
	m1 := make(map[string]interface{})
	m1["jj"] = "m1jjj"
	m1["dd"] = "m1ddd"
	m2 := make(map[string]interface{})
	m2["jj"] = "m2jjj"
	m2["dd"] = "m2ddd"
	s = append(s, m1)
	s = append(s, m2)
	//测试自定义异常处理函数
	//throw.Handle = func(p *parameter.Parameter) {
	//	if p.Name == throw.MESSAGE {
	//		throw.LogPayload.Msg = throw.ErrorToString(p.Value)
	//		throw.MapString.Put(p.Name, p.Value)
	//	}
	//}
	Exception.BusinessException(exception.Code(2), exception.Message("v3 https 测试异常2"))
	exception.Throw(exception.Message("v3 https 测试异常11"), exception.Code(2), exception.Data(struct {
		Id   int
		Name string
	}{Id: 1, Name: "哦哦"}))
	return "v3 https_test_throw"
}

//错误码
type ErrorCode int64

func (this ErrorCode) Message(variables ...interface{}) string {
	return code.Get(this, variables...)
}

func (this ErrorCode) Register() code.Message {
	code400001()
	return code.Container()
}

const (
	MobileEmpty   ErrorCode = iota + 400001 //mobile错误
	PasswordError                           //password错误
	UseridsError                            //user_ids错误
	NameError                               //name错误
	MinError                                //不能小于4位
)

func code400001() {
	code.Container().
		Put(MobileEmpty, "mobile错误").
		Put(PasswordError, "password错误").
		Put(UseridsError, "user_ids错误").
		Put(NameError, "name错误").
		Put(MinError, "不能小于4位")
}

//测试验证器
type DutyUser struct {
	DutyUserId       int64   `json:"duty_user_id" binding:"mobile"`
	EducationClassId int64   `json:"education_class_id" binding:"required"`
	UserIds          []int64 `json:"user_ids" binding:"user_ids"`
	User             User    `json:"user"`
}

type User struct {
	Name string `json:"name" binding:"name"`
}

func NewDutyUser() *DutyUser {
	return &DutyUser{UserIds: make([]int64, 0)}
}

func (this *DutyUser) RegisterValidator() *higo.Verify {
	return higo.Verifier().
		Tag("mobile",
			higo.Rule("required", MobileEmpty)).
		Tag("password",
			higo.Rule("min=4", func() higo.ValidatorToFunc {
				return func(fl validator.FieldLevel) (bool, code.ICode) {
					fmt.Println("DemoController:123", fl.Field().Interface())
					return true, MinError
				}
			}()),
			higo.Rule("required", PasswordError),
		).
		Tag("user_ids",
			higo.Rule("required", UseridsError)).
		Tag("name",
			higo.RuleFunc("required", func(fl validator.FieldLevel) (b bool, iCode code.ICode) {
				fmt.Println("DemoController:140", fl.Field().Interface())
				return true, NameError
			}))
}

// 测试校验
func (this *DemoController) HttpsTestValidate() {
	ctx := request.Context()
	param := NewDutyUser()
	//param.DutyUserId = 1
	//param.EducationClassId = 2
	fmt.Println("DemoController:139", higo.VerifyContainer)
	//校验数据
	//higo.Receiver(ctx.ShouldBindJSON(param)).Unwrap()
	//higo.Validate(param).Unwrap()
	higo.Validate(param).Receiver(ctx.ShouldBindJSON(param)).Unwrap()

	log.Fatalln(param)
}

// 测试get请求
func (this *DemoController) HttpsTestGet(ctx *gin.Context) higo.Model {
	param := NewDutyUser()
	fmt.Println(higo.VerifyContainer)
	//校验数据
	//higo.Receiver(ctx.ShouldBindJSON(param)).Unwrap()
	higo.Validate(param).Unwrap()
	//higo.Validate(param).Receiver(ctx.ShouldBindJSON(param)).Unwrap()

	log.Fatalln(param)
	/**
	fmt.Println(injector.BeanFactory.Get(this))
	fmt.Println(this)
	fmt.Printf("%p\n", this)
	*/
	user := UserModel.New(UserModel.WithId(101))
	/**
	user.Uname = this.Age.String()
	fmt.Println(user)
	higo.Result(ctx.ShouldBindJSON(user)).Unwrap()
	err := ctx.ShouldBindUri(user)
	if err != nil {
		log.Fatal("映射错误")
	}
	user.UserById(3, "*")
	fmt.Println(user)
	*/
	user.Add("werwerwerg123456", "15987", 20)
	//this.Table("ts_user").
	//	Where("id=?", 3).
	//	Find(user)
	higo.Task(this.TestTask, func() {
		this.TestTaskDone(3)
	}, user.Id)
	//redisConn := this.Pool.Get()
	//fmt.Println(redis.String(redisConn.Do("get", "name")))
	return user
}

//测试发布
func (this *DemoController) TestPub(ctx *gin.Context) string {
	return "测试发布"
}

// 测试post请求
func (this *DemoController) HttpsTestPost(ctx *gin.Context) string {
	fmt.Println(this)
	this.Name = "1001"
	fmt.Println(this)
	fmt.Printf("%p\n", this)
	return "v3 https_test_post"
}

// 测试异常
func (this *DemoController) HttpTestThrow(ctx *gin.Context) string {
	exception.Throw(exception.Message("v3 http 测试异常"), exception.Code(0))
	return "v3 http_test_throw"
}

// 测试get请求
func (this *DemoController) HttpTestGet(ctx *gin.Context) string {
	return "HttpTestGet"
}

// 测试post请求
func (this *DemoController) HttpTestPost(ctx *gin.Context) string {
	return "v3 http_test_post"
}

func (this *DemoController) Login(ctx *gin.Context) string {
	fmt.Println(this)
	this.Name = "1000"
	fmt.Println(this)
	fmt.Printf("%p\n", this)
	return "登录成功11"
}

func (this *DemoController) Login1(ctx *gin.Context) {
	fmt.Println(this)
	this.Name = "1000"
	fmt.Println(this)
	fmt.Printf("%p\n", this)
	//higo.Responser(ctx).SuccessJson(this.Name, 10000, nil)
	fmt.Println(11)
	higo.Responser(ctx).ErrorJson(this.Name, 10000, nil)
}

func (this *DemoController) TestTask(params ...interface{}) {
	time.Sleep(time.Second * 5)
	fmt.Println("测试task")
	fmt.Println(params)
}

func (this *DemoController) TestTaskDone(id int) {
	fmt.Println("测试task执行结束", id)
}
