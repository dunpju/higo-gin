package higo

import (
	"fmt"
	"github.com/dengpju/higo-config/config"
	iocConfig "github.com/dengpju/higo-ioc/config"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/dengpju/higo-logger/logger"
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-throw/exception"
	"github.com/dengpju/higo-utils/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var (
	hg *Higo
	// ssl 证书
	SslOut, SslCrt, SslKey string
	isLoadEnv              bool
)

type Higo struct {
	*gin.Engine
	group       *gin.RouterGroup
	errgroup    errgroup.Group
	isAutoTLS   bool
	bits        int
	isRedisPool bool
	middle      []IMiddleware
	serveType   string // serve type
	serve       string // serve name
}

// 初始化
func Init(root *utils.SliceString) *Higo {
	hg = &Higo{
		Engine: gin.New(),
		middle: make([]IMiddleware, 0),
		serve:  router.DefaultServe,
		bits:   1024,
	}

	// 全局异常
	hg.Engine.Use(NewRecover().Exception(hg))
	//设置跨域、鉴权
	hg.Middleware(NewCors(), NewAuth())
	// 初始分隔符
	pathSeparator = utils.PathSeparator()
	AppConfigDir.Clone(root)
	root.ForEach(func(index int, value interface{}) {
		AppConfigDir.Append(value)
	})
	AppConfigDir.Append("app")
	AppConfigDir.Append("Config")
	// 是否使用自带ssl测试https
	hg.isAutoTLS = false
	// 未加载env
	if false == isLoadEnv {
		hg.LoadEnv(root)
	}

	return hg
}

func (this *Higo) SetPathSeparator(sep string) *Higo {
	pathSeparator = sep
	return this
}

func (this *Higo) PathSeparator() string {
	return pathSeparator
}

func (this *Higo) SetBits(bits int) *Higo {
	this.bits = bits
	return this
}

// 设置主目录
func (this *Higo) setRoot(r *utils.SliceString) *Higo {
	root = r
	return this
}

// 获取主目录
func (this *Higo) GetRoot() *utils.SliceString {
	return Root()
}

// 加载env
func (this *Higo) LoadEnv(root *utils.SliceString) *Higo {

	utils.SetPathSeparator(pathSeparator)
	// 设置主目录
	this.setRoot(root)
	// 创建runtime
	utils.Mkdir(this.GetRoot().Separator(pathSeparator) + "runtime")
	// 日志
	logger.Logrus.Root(this.GetRoot().Separator(pathSeparator)).File("higo").Init()
	// 装载env配置
	env := this.GetRoot().Separator(pathSeparator) + "env"
	fmt.Println(env)
	if ! utils.DirExist(env) {
		utils.Mkdir(env)
	}
	envConf := config.New()
	filepathErr := filepath.Walk(env,
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			if path.Ext(p) == ".yaml" {
				yamlFile, err := ioutil.ReadFile(p)
				if err != nil {
					logger.LoggerStack(err, utils.GoroutineID())
				}
				yamlMap := make(map[interface{}]interface{})
				yamlFileErr := yaml.Unmarshal(yamlFile, yamlMap)
				envConf.Set(utils.Basename(p, "yaml"), yamlMap)
				if yamlFileErr != nil {
					exception.Throw(exception.Message(yamlFileErr), exception.Code(0))
				}
				logger.Logrus.Infoln("Loader env file:", filepath.Base(p))
			}
			return nil
		})
	if filepathErr != nil {
		exception.Throw(exception.Message(filepathErr), exception.Code(0))
	}

	config.Set(config.EnvConf, envConf)
	this.loadConfigur()
	SslOut = this.GetRoot().Separator(pathSeparator) + config.App("SSL.OUT").(string) + pathSeparator
	SslCrt = config.App("SSL.CRT").(string)
	SslKey = config.App("SSL.KEY").(string)

	// bean
	this.Beans(NewBean())

	isLoadEnv = true

	return this
}

// 加载配置
func (this *Higo) loadConfigur() *Higo {
	if ! utils.DirExist(AppConfigDir.Separator(utils.PathSeparator())) {
		utils.Mkdir(AppConfigDir.Separator(utils.PathSeparator()))
	}
	conf := config.New()
	filepathErr := filepath.Walk(AppConfigDir.Separator(utils.PathSeparator()),
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			if path.Ext(p) == ".yaml" {
				yamlFile, err := ioutil.ReadFile(p)
				if err != nil {
					logger.LoggerStack(err, utils.GoroutineID())
				}
				yamlMap := make(map[interface{}]interface{})
				yamlFileErr := yaml.Unmarshal(yamlFile, yamlMap)
				conf.Set(utils.Basename(p, "yaml"), yamlMap)
				if yamlFileErr != nil {
					exception.Throw(exception.Message(yamlFileErr), exception.Code(0))
				}
				logger.Logrus.Infoln("Loader config file:", filepath.Base(p))
			}
			return nil
		})
	if filepathErr != nil {
		exception.Throw(exception.Message(filepathErr), exception.Code(0))
	}
	config.Set("config", conf)
	return this
}

//全局中间件
func (this *Higo) Middleware(middlewares ...IMiddleware) *Higo {
	this.middle = append(this.middle, middlewares...)
	return this
}

func (this *Higo) SetName(serve string) *Higo {
	this.serve = serve
	return this
}

// 获取serve name
func (this *Higo) Serve() string {
	return this.serve
}

func (this *Higo) Type() string {
	return this.serveType
}

func (this *Higo) SetType(serveType string) *Higo {
	this.serveType = serveType
	return this
}

func (this *Higo) AddServe(route IRouterLoader, middles ...IMiddleware) *Higo {
	injector.BeanFactory.Apply(route)
	injector.BeanFactory.Set(route)
	if ! onlySupportServe.Exist(route.GetServe().Type) {
		panic("Serve Type error! only support:" + onlySupportServe.String() + ", But give " + route.GetServe().Type)
	}
	route.GetServe().SetRouter(route).SetMiddle(middles...)
	serves = append(serves, route.GetServe())
	return this
}

// 是否自动生成ssl证书
func (this *Higo) IsAutoTLS(isAuto bool) *Higo {
	this.isAutoTLS = isAuto
	return this
}

// 使用redis连接池
func (this *Higo) IsRedisPool() *Higo {
	this.isRedisPool = true
	return this
}

//启动
func (this *Higo) Boot() {
	//启动服务
	for _, ser := range serves {
		//初始化
		hg := Init(this.GetRoot()).
			//设置服务类型
			SetType(ser.Type).
			//设置服务名称
			SetName(ser.Name)

		//全局中间件
		for _, m := range this.middle {
			hg.Engine.Use(m.Middle(hg))
		}
		//服务中间件
		for _, m := range ser.Middle {
			hg.Engine.Use(m.Middle(hg))
		}
		//是否使用自带ssl测试https
		if this.isAutoTLS {
			//生成ssl证书
			utils.NewTLS(SslOut, SslCrt, SslKey).SetBits(this.bits).Build()
		}
		//是否使用redis pool
		if this.isRedisPool {
			InitRedisPool()
		}
		//运行模式debug/release
		if gin.ReleaseMode == config.App("MODE") {
			gin.SetMode(gin.ReleaseMode)
		}

		configs := config.Get(ser.Config).(*config.Configure)
		addr := configs.Get("Addr").(string)
		readTimeout := configs.Get("ReadTimeout").(int)
		writeTimeout := configs.Get("WriteTimeout").(int)

		//添加服务
		router.AddServe(hg.serve)
		handler := ser.Router.Loader(hg)
		//加载路由
		hg.loadRoute()

		serve := &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
		}

		if ser.Type == HttpServe {
			this.errgroup.Go(func() error {
				logger.Logrus.Infoln("HTTP Server listening at " + addr + " Starting Success!")
				return serve.ListenAndServe()
			})
		}
		if ser.Type == HttpsServe {
			this.errgroup.Go(func() error {
				logger.Logrus.Infoln("HTTPS Server listening at " + addr + " Starting Success!")
				return serve.ListenAndServeTLS(SslOut+SslCrt, SslOut+SslKey)
			})
		}
		if ser.Type == WebsocketServe {
			this.errgroup.Go(func() error {
				logger.Logrus.Infoln("WEBSOCKET Server listening at " + addr + " Starting Success!")
				return serve.ListenAndServe()
			})
		}
	}

	//启动定时任务
	CronTask().Start()

	if err := this.errgroup.Wait(); err != nil {
		logger.Logrus.Fatal(err)
	}
}

// 获取路由
func (this *Higo) GetRoute(relativePath string) (*router.Route, bool) {
	return RouterContainer.Get(relativePath), true
}

// 静态文件
func (this *Higo) StaticFile(relativePath, filepath string) *Higo {
	// 添加路由容器
	router.AddRoute(router.GET, relativePath, "", router.IsStatic(true), router.SetServe(this.serve))
	hg.Engine.StaticFile(relativePath, filepath)
	return this
}

// 装载路由
func (this *Higo) loadRoute() *Higo {
	router.GetRoutes(this.serve).ForEach(func(index int, route *router.Route) {
		// 判断空标记
		IsEmptyFlag(route)
		// 添加路由容器
		RouterContainer.Add(route.Prefix()+route.RelativePath(), route)
		if route.Prefix() != "" {
			this.group = this.Engine.Group(route.Prefix())
			this.GroupHandle(route)
		} else {
			this.Handle(route)
		}
	})
	return this
}

//注入
func (this *Higo) Inject(classs ...IClass) *Higo {
	for _, class := range classs {
		injector.BeanFactory.Apply(class)
		injector.BeanFactory.Set(class)
		AddContainer(class)
	}
	return this
}

func (this *Higo) Route(controllers ...IController) *Higo {
	for _, controller := range controllers {
		injector.BeanFactory.Apply(controller)
		injector.BeanFactory.Set(controller)
		AddContainer(controller)
		controller.Route(hg)
	}
	return this
}

// 路由组Handle
func (this *Higo) GroupHandle(route *router.Route) *Higo {
	if handle := Convert(route.Handle()); handle != nil {
		handles := appendHandle(handle, route)
		this.group.Handle(strings.ToUpper(route.Method()), route.RelativePath(), handles...)
	}
	return this
}

// 路由Handle
func (this *Higo) Handle(route *router.Route) *Higo {
	if handle := Convert(route.Handle()); handle != nil {
		handles := appendHandle(handle, route)
		this.Engine.Handle(strings.ToUpper(route.Method()), route.RelativePath(), handles...)
	}
	return this
}

func appendHandle(handle gin.HandlerFunc, route *router.Route) []gin.HandlerFunc {
	handles := handleSlice(route)
	if reflect.ValueOf(route.Handle()).Type().ConvertibleTo(refWsResponder) {
		handles = append(handles, wsUpgraderHandle(route))
	} else {
		handles = append(handles, handle)
	}
	return handles
}

func handleSlice(route *router.Route) []gin.HandlerFunc {
	handles := make([]gin.HandlerFunc, 0)
	if reflect.ValueOf(route.Handle()).Type().ConvertibleTo(refWsResponder) {
		handles = append(handles, wsConnMiddleWare())
	}
	if groupMiddles, ok := route.GroupMiddle().([]interface{}); ok {
		for _, groupMiddle := range groupMiddles {
			if middle, ok := groupMiddle.(gin.HandlerFunc); ok {
				handles = append(handles, middle)
			}
		}
	}
	if middlewares, ok := route.Middleware().([]interface{}); ok {
		for _, middleware := range middlewares {
			if middle, ok := middleware.(gin.HandlerFunc); ok {
				handles = append(handles, middle)
			}
		}
	}
	return handles
}

// 添加到Bean
func (this *Higo) Beans(configs ...iocConfig.IBean) *Higo {
	for _, conf := range configs {
		injector.BeanFactory.Config(conf)
	}
	return this
}

// 定时任务
func (this *Higo) Cron(expr string, fn func()) *Higo {
	_, err := CronTask().AddFunc(expr, fn)
	if err != nil {
		exception.Throw(exception.Message(err), exception.Code(0))
	}
	return this
}
