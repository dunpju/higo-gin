package higo

import (
	"fmt"
	"github.com/dunpju/higo-annotation/anno"
	"github.com/dunpju/higo-config/config"
	"github.com/dunpju/higo-ioc/injector"
	"github.com/dunpju/higo-logger/logger"
	"github.com/dunpju/higo-router/router"
	"github.com/dunpju/higo-throw/exception"
	"github.com/dunpju/higo-utils/utils"
	"github.com/dunpju/higo-utils/utils/dirutil"
	"github.com/dunpju/higo-utils/utils/runtimeutil"
	"github.com/dunpju/higo-utils/utils/sliceutil"
	"github.com/dunpju/higo-utils/utils/tlsutil"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

var (
	hg                     *Higo
	SslOut, SslCrt, SslKey string // ssl
	isLoadEnv              bool
	isLoadAppConfig        bool
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
	serves      []serve
	root        *sliceutil.SliceString
}

type serve struct {
	router  IRouterLoader
	middles []IMiddleware
}

func _newServe(router IRouterLoader, middles []IMiddleware) serve {
	return serve{router: router, middles: middles}
}

// Init 初始化
func Init(root *sliceutil.SliceString) *Higo {
	hg = &Higo{
		Engine: gin.New(),
		middle: make([]IMiddleware, 0),
		serves: make([]serve, 0),
		serve:  router.DefaultServe,
		bits:   1024,
		root:   root,
	}
	return hg
}

func (this *Higo) _init() *Higo {
	// 全局异常
	this.Engine.Use(NewRecover().Exception(this))
	// 设置跨域、鉴权
	this.Middleware(NewCors(), NewAuth())
	// 初始分隔符
	pathSeparator = dirutil.PathSeparator()
	// 未加载env
	if !isLoadEnv {
		this.LoadEnv(this.root)
		isLoadEnv = true
	}
	if !isLoadAppConfig {
		AppConfigDir.Clone(this.root)
		this.root.ForEach(func(index int, value interface{}) bool {
			AppConfigDir.Append(value.(string))
			return true
		})
		appConfig := config.Env("app.APP_CONFIG")
		if nil == appConfig {
			AppConfigDir.Append("app")
			AppConfigDir.Append("config")
		} else {
			AppConfigDir.Append(appConfig.(string))
		}

		isLoadAppConfig = true

		this.loadConfigure()
		eventPoint(this, AfterLoadConfigure)

		SslOut = this.GetRoot().Join(pathSeparator) + config.App("SSL.OUT").(string) + pathSeparator
		SslCrt = config.App("SSL.CRT").(string)
		SslKey = config.App("SSL.KEY").(string)
		// higo bean
		this.Beans(NewBean())
	}
	return this
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
func (this *Higo) setRoot(r *sliceutil.SliceString) *Higo {
	root = r
	return this
}

// GetRoot 获取主目录
func (this *Higo) GetRoot() *sliceutil.SliceString {
	return Root()
}

type yamlRaw struct {
	parent         *yamlRaw
	prefixBlankNum int
	key            string
	value          interface{}
	child          []*yamlRaw
}

// LoadEnv 加载env
func (this *Higo) LoadEnv(root *sliceutil.SliceString) *Higo {
	dirutil.SetPathSeparator(pathSeparator)
	// 设置主目录
	this.setRoot(root)
	// 创建runtime
	dirutil.Mkdir(this.GetRoot().Join(pathSeparator) + "runtime")
	// 日志
	logger.Logrus.Root(this.GetRoot().Join(pathSeparator)).File("higo").Init()
	// 装载env配置
	envPath := this.GetRoot().Join(pathSeparator) + "env"
	if !dirutil.DirExist(envPath) {
		dirutil.Mkdir(envPath)
	}
	envConf := config.New()
	filepathErr := filepath.Walk(envPath,
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			if path.Ext(p) == ".yaml" {
				fileBase := filepath.Base(p)
				ok := strings.HasSuffix(envPath+pathSeparator+fileBase, p)
				if ok {
					yamlFile, err := ioutil.ReadFile(p)
					if err != nil {
						goid, _ := runtimeutil.GoroutineID()
						logger.LoggerStack(err, goid)
					}
					yamlMap := make(map[interface{}]interface{})
					yamlFileErr := yaml.Unmarshal(yamlFile, yamlMap)
					envConf.Set(utils.Dir.Basename(p, "yaml"), yamlMap)
					if yamlFileErr != nil {
						exception.Throw(exception.Message(yamlFileErr), exception.Code(0))
					}
					logger.Logrus.Infoln("Loader env config file:", p)
				}
			}
			return nil
		})
	if filepathErr != nil {
		exception.Throw(exception.Message(filepathErr), exception.Code(0))
	}

	config.Set(config.EnvConf, envConf)
	return this
}

// 加载配置
func (this *Higo) loadConfigure() *Higo {
	if !dirutil.DirExist(AppConfigDir.Join(dirutil.PathSeparator())) {
		dirutil.Mkdir(AppConfigDir.Join(dirutil.PathSeparator()))
	}
	conf := config.New()
	filepathErr := filepath.Walk(AppConfigDir.Join(dirutil.PathSeparator()),
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			if path.Ext(p) == ".yaml" {
				yamlFile, err := os.ReadFile(p)
				if err != nil {
					goid, _ := runtimeutil.GoroutineID()
					logger.LoggerStack(err, goid)
				}
				yamlMap := make(map[interface{}]interface{})
				yamlFileErr := yaml.Unmarshal(yamlFile, yamlMap)
				conf.Set(dirutil.Basename(p, "yaml"), yamlMap)
				if yamlFileErr != nil {
					exception.Throw(exception.Message(yamlFileErr), exception.Code(0))
				}
				logger.Logrus.Infoln("Loader app config file:", p)
			}
			return nil
		})
	if filepathErr != nil {
		exception.Throw(exception.Message(filepathErr), exception.Code(0))
	}
	config.Set("config", conf)
	anno.Config = config.Anno("value").(*config.Configure)
	return this
}

// Middleware 全局中间件
func (this *Higo) Middleware(middlewares ...IMiddleware) *Higo {
	for _, middleware := range middlewares {
		if m, ok := middleware.(*Cors); ok && len(this.middle) > 0 {
			middle := make([]IMiddleware, 0)
			middle = append(middle, m)
			this.middle = append(middle, this.middle...)
		} else {
			this.middle = append(this.middle, middleware)
		}
	}
	return this
}

// AuthHandlerFunc 设置鉴权中间件
func (this *Higo) AuthHandlerFunc(middle IMiddleware) *Higo {
	MiddleAuthFunc = middle.Middle
	return this
}

// CorsHandlerFunc 设置跨域中间件
func (this *Higo) CorsHandlerFunc(middle IMiddleware) *Higo {
	MiddleCorsFunc = middle.Middle
	return this
}

// SetName 设置服务名称
func (this *Higo) SetName(serve string) *Higo {
	this.serve = serve
	return this
}

// Serve 获取serve name
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
	this.serves = append(this.serves, _newServe(route, middles))
	return this
}

func (this *Higo) registerServe(route IRouterLoader, middles ...IMiddleware) *Higo {
	injector.BeanFactory.Apply(route)
	injector.BeanFactory.Set(route)
	if !onlySupportServe.Exist(route.GetServe().Type) {
		panic("Serve Type error! only support:" + onlySupportServe.String() + ", But give " + route.GetServe().Type)
	}
	route.GetServe().SetRouter(route).SetMiddle(middles...)
	serves = append(serves, route.GetServe())
	return this
}

// IsAutoTLS 是否自动生成ssl证书
func (this *Higo) IsAutoTLS(isAuto bool) *Higo {
	this.isAutoTLS = isAuto
	return this
}

// IsRedisPool 使用redis连接池
func (this *Higo) IsRedisPool() *Higo {
	this.isRedisPool = true
	return this
}

func (this *Higo) InitGroupIsAuth(b bool) *Higo {
	router.SetInitGroupIsAuth(b)
	return this
}

// Boot 启动
func (this *Higo) Boot() {
	// 初始化
	this._init()
	// 注册服务
	for _, s := range this.serves {
		this.registerServe(s.router, s.middles...)
	}

	//自动注册校验
	VerifyContainer.Range(func(key, verify interface{}) bool {
		verify.(*Verify).VerifyRules.Range(func(tag, rules interface{}) bool {
			RegisterValidation(tag.(string), rules.(*RuleGroup).ToFunc())
			return true
		})
		return true
	})
	//启动服务
	for _, ser := range serves {
		//初始化
		hg := Init(this.GetRoot()).
			_init().
			//设置服务类型
			SetType(ser.Type).
			//设置服务名称
			SetName(ser.Name)

		//全局中间件
		for _, m := range this.middle {
			hg.Engine.Use(m.Middle(hg))
		}
		//serve Middle
		for _, mid := range ser.Middle {
			if m, ok := mid.(IMiddleware); ok {
				hg.Engine.Use(m.Middle(hg))
			}
		}
		//是否使用自带ssl测试https
		if this.isAutoTLS {
			//生成ssl证书
			tlsutil.NewTLS(SslOut, SslCrt, SslKey).SetBits(this.bits).Build()
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

		//添加服务
		router.AddServe(hg.serve)
		//serve bean router
		for _, mid := range ser.Middle {
			if bean, ok := mid.(injector.IBean); ok {
				hg.Beans(bean)
			}
		}

		this.run(func() {
			eventPoint(hg, BeforeLoadRoute)
			ser.Router.Loader(hg)
			hg.loadRoute()
			eventPoint(hg, AfterLoadRoute)
			if ser.Type == HttpServe {
				this.errgroup.Go(func() error {
					logger.Logrus.Infoln("HTTP Server listening at " + addr + " Starting Success!")
					return hg.Run(addr)
				})
			}
			if ser.Type == HttpsServe {
				this.errgroup.Go(func() error {
					logger.Logrus.Infoln("HTTPS Server listening at " + addr + " Starting Success!")
					return hg.RunTLS(addr, SslOut+SslCrt, SslOut+SslKey)
				})
			}
			if ser.Type == WebsocketServe {
				this.errgroup.Go(func() error {
					logger.Logrus.Infoln("WEBSOCKET Server listening at " + addr + " Starting Success!")
					return hg.Run(addr)
				})
			}
		})
	}

	//启动定时任务
	CronTask().Start()

	if err := this.errgroup.Wait(); err != nil {
		logger.Logrus.Fatal(err)
	}
}

func (this *Higo) run(fn func()) {
	if len(os.Args) <= 1 {
		fn()
	} else {
		match, err := regexp.Match(`^[(gen)|(-g)]`, []byte(os.Args[1]))
		if err != nil {
			fmt.Println(err)
		}
		if match {
			RootInit()
			if err := RootCommand.Execute(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Exit(1)
		}
	}
}

func (this *Higo) Event(eventType EventType, handle EventHandle) *Higo {
	addEvent(eventType, handle)
	return this
}

// GetRoute 获取路由
func (this *Higo) GetRoute(method, relativePath string) (*router.Route, bool) {
	return RouterContainer.Get(this.serve, method, relativePath), true
}

// StaticFile 静态文件
func (this *Higo) StaticFile(relativePath, filepath string) *Higo {
	// 添加路由容器
	router.AddRoute(router.GET, relativePath, "", router.IsStatic(true), router.SetServe(this.serve))
	hg.Engine.StaticFile(relativePath, filepath)
	return this
}

// Static 静态目录
func (this *Higo) Static(relativePath, root string) *Higo {
	// 添加路由容器
	router.AddRoute(router.GET, relativePath, "", router.IsStatic(true), router.SetServe(this.serve))
	hg.Engine.Static(relativePath, root)
	return this
}

// StaticFS 静态目录
func (this *Higo) StaticFS(relativePath string, fs http.FileSystem) *Higo {
	// 添加路由容器
	router.AddRoute(router.GET, relativePath, "", router.IsStatic(true), router.SetServe(this.serve))
	hg.Engine.StaticFS(relativePath, fs)
	return this
}

// 装载路由
func (this *Higo) loadRoute() *Higo {
	router.GetRoutes(this.serve).ForEach(func(route *router.Route) {
		// 判断空标记
		IsEmptyFlag(route)

		if route.Prefix() != "" {
			this.group = this.Engine.Group(route.Prefix())
			this.GroupHandle(route)
		} else {
			this.Handle(route)
		}
	})
	return this
}

// Register register to di container
func Register(classs ...IClass) {
	for _, class := range classs {
		AddContainer(class.New)
	}
}

func (this *Higo) register(conf injector.IBean) *Higo {
	t := reflect.TypeOf(conf)
	if t.Kind() != reflect.Ptr {
		panic("required ptr object")
	}
	v := reflect.ValueOf(conf)
	for i := 0; i < t.NumMethod(); i++ {
		method := v.Method(i)
		typeRegexp := regexp.MustCompile(`func\((.*)\)`)
		regParams := typeRegexp.FindStringSubmatch(fmt.Sprintf("%s", method.Type()))
		if "" != regParams[1] { // 有参数
			arguments := make([]reflect.Value, 0)
			args := strings.Split(regParams[1], ",")
			for _, a := range args {
				trimArgType := strings.Trim(a, " ")
				if "string" == trimArgType {
					arguments = append(arguments, reflect.ValueOf(""))
				} else if "int" == trimArgType {
					arguments = append(arguments, reflect.ValueOf(0))
				} else if "int64" == trimArgType {
					arguments = append(arguments, reflect.ValueOf(int64(0)))
				}
			}
			if len(arguments) > 0 {
				callRet := method.Call(arguments)
				if callRet != nil && len(callRet) == 1 {
					if class, ok := callRet[0].Interface().(IClass); ok {
						Register(class)
					}
					if controller, ok := callRet[0].Interface().(IController); ok {
						this.Route(controller)
					}
				}
			}
		} else { // 无参数
			callRet := method.Call(nil)
			if callRet != nil && len(callRet) == 1 {
				if class, ok := callRet[0].Interface().(IClass); ok {
					Register(class)
				}
				if controller, ok := callRet[0].Interface().(IController); ok {
					this.Route(controller)
				}
			}
		}
	}
	return this
}

func (this *Higo) Route(controllers ...IController) *Higo {
	for _, controller := range controllers {
		AddContainer(controller.New)
		injector.BeanFactory.Apply(controller)
		injector.BeanFactory.Set(controller)
		controller.Route(this)
	}
	return this
}

// GroupHandle 路由组Handle
func (this *Higo) GroupHandle(route *router.Route) *Higo {
	if handle := Convert(route.Handle()); handle != nil {
		handles := appendHandle(handle, route)
		this.group.Handle(strings.ToUpper(route.Method()), route.RelativePath(), handles...)
	}
	return this
}

// Handle 路由Handle
func (this *Higo) Handle(route *router.Route) *Higo {
	if handle := Convert(route.Handle()); handle != nil {
		handles := appendHandle(handle, route)
		this.Engine.Handle(strings.ToUpper(route.Method()), route.RelativePath(), handles...)
	}
	return this
}

// 中间件顺序倒序包裹，越往后添加的中间件越贴近需要执行的逻辑
func appendHandle(handle gin.HandlerFunc, route *router.Route) []gin.HandlerFunc {
	handles := handleSlice(route)
	if reflect.ValueOf(route.Handle()).Type().ConvertibleTo(refWsResponder) {
		handles = append(handles, wsUpgraderHandle(route))
	} else {
		handles = append(handles, handle)
	}
	return handles
}

// handle切片
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

// Beans 添加Bean
func (this *Higo) Beans(configs ...injector.IBean) *Higo {
	for _, conf := range configs {
		injector.BeanFactory.Config(conf)
		this.register(conf)
	}
	return this
}

// Cron 定时任务
func (this *Higo) Cron(expr string, fn func()) *Higo {
	_, err := CronTask().AddFunc(expr, fn)
	if err != nil {
		exception.Throw(exception.Message(err), exception.Code(0))
	}
	return this
}
