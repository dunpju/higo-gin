package higo

import (
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
func Init() *Higo {
	hg = &Higo{
		Engine: gin.New(),
		middle: make([]IMiddleware, 0),
		serve:  router.DefaultServe,
		bits:   1024,
	}

	// 全局异常
	hg.Engine.Use(NewRecover().Exception(hg))
	// 初始分隔符
	pathSeparator = utils.PathSeparator()
	// 是否使用自带ssl测试https
	hg.isAutoTLS = false

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
				yamlFile, _ := ioutil.ReadFile(p)
				conf := config.New()
				yamlFileErr := yaml.Unmarshal(yamlFile, conf)
				envConf.Set(utils.Basename(p, "yaml"), conf)
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
	config.Set("env", envConf)
	SslOut = this.GetRoot().Separator(pathSeparator) + config.String("env.app.SSL.OUT") + pathSeparator
	SslCrt = config.String("env.app.SSL.CRT")
	SslKey = config.String("env.app.SSL.KEY")
	return this
}

// 加载配置
func (this *Higo) LoadConfigur(root *utils.SliceString) *Higo {
	return this
}

// 中间件装载器
func (this *Higo) Middleware(imiddleware ...IMiddleware) *Higo {
	for _, middleware := range imiddleware {
		this.middle = append(this.middle, middleware)
	}
	return this
}

func (this *Higo) SetServe(serve string) *Higo {
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

func (this *Higo) AddServe(router IRouterLoader) *Higo {
	if ! onlySupportServe.Exist(router.Serve().Type) {
		panic("Serve Type error! only support:" + onlySupportServe.String() + ", But give " + router.Serve().Type)
	}
	serves = append(serves, router.Serve())
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

// 启动
func (this *Higo) Boot() {
	// 启动服务
	for _, ser := range serves {
		logger.Logrus.Infoln("Server starting......")
		// 初始化
		hg := Init().
			//设置服务类型
			SetType(ser.Type).
			//设置服务名称
			SetServe(ser.Name).
			//加载配置
			LoadConfigur(this.GetRoot())
		// 中间件
		for _, m := range this.middle {
			hg.Engine.Use(m.Loader(hg))
		}
		// 是否使用自带ssl测试https
		if this.isAutoTLS {
			// 生成ssl证书
			utils.NewTLS(SslOut, SslCrt, SslKey).SetBits(this.bits).Build()
		}
		// 是否使用redis pool
		if this.isRedisPool {
			InitRedisPool()
		}
		// 运行模式debug/release
		if gin.ReleaseMode == config.String("env.app.MODE") {
			gin.SetMode(gin.ReleaseMode)
		}

		configs := config.Get(ser.Config).(config.Configure)
		addr := configs.Get("Addr").(string)
		readTimeout := configs.Get("ReadTimeout").(int)
		writeTimeout := configs.Get("WriteTimeout").(int)

		// 添加服务
		router.AddServe(hg.serve)
		handler := ser.Router.Loader(hg)
		// 加载路由
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

	// 启动定时任务
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
