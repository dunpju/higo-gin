package higo

import (
	"fmt"
	"gitee.com/dengpju/higo-configure/configure"
	"github.com/dengpju/higo-gin/higo/consts"
	iocConfig "github.com/dengpju/higo-ioc/config"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/dengpju/higo-logger/logger"
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-throw/throw"
	"github.com/dengpju/higo-utils/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	hg *Higo
	// 路径分隔符
	PathSeparator string
	// ssl 证书
	SslOut, SslCrt, SslKey string
	Once                   sync.Once
)

// http 服务结构体
type Hse struct {
	Config string
	Router IRouterLoader
	Serve  string
}

type Higo struct {
	*gin.Engine
	group       *gin.RouterGroup
	errgroup    errgroup.Group
	root        string
	isAutoTLS   bool
	bits        int
	isRedisPool bool
	middle      []IMiddleware
	serve       []Hse
}

// 初始化
func Init() *Higo {
	hg = &Higo{
		Engine: gin.New(),
		middle: make([]IMiddleware, 0),
		serve:  make([]Hse, 0),
		bits:   1024,
	}

	// 全局异常
	hg.Engine.Use(NewRecover().Exception(hg))
	// 初始分隔符
	PathSeparator = string(os.PathSeparator)
	// 是否使用自带ssl测试https
	hg.isAutoTLS = false

	return hg
}

func (this *Higo) SetBits(bits int) *Higo {
	this.bits = bits
	return this
}

// 设置主目录
func (this *Higo) setRoot(root string) *Higo {
	this.root = root
	return this
}

// 获取主目录
func (this *Higo) GetRoot() string {
	return utils.If(this.root == "", consts.ROOT, this.root).(string)
}

// 加载配置
func (this *Higo) LoadConfigur(root string) *Higo {
	// 设置主目录
	this.setRoot(root)
	// runtime目录
	runtimeDir := root + "runtime"
	if _, err := os.Stat(runtimeDir); os.IsNotExist(err) {
		if os.Mkdir(runtimeDir, os.ModePerm) != nil {
		}
	}
	// 日志
	logger.Logrus.Root(root).File("higo").Init()
	// 装载env配置
	confDir := root + "env"
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		if err = os.Mkdir(confDir, os.ModePerm); err != nil {
			throw.Throw(throw.Message(err), throw.Code(0))
		}
	}
	filepathErr := filepath.Walk(confDir,
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			if path.Ext(p) == ".yaml" {
				logger.Logrus.Infoln("Loader Configure file:", filepath.Base(p))
				yamlFile, _ := ioutil.ReadFile(p)
				yamlFileErr := yaml.Unmarshal(yamlFile, configure.New())
				if yamlFileErr != nil {
					throw.Throw(throw.Message(yamlFileErr), throw.Code(0))
				}
			}
			return nil
		})
	if filepathErr != nil {
		throw.Throw(throw.Message(filepathErr), throw.Code(0))
	}
	mapSslConf := configure.Config("SSL")
	SslOut = root + mapSslConf.String("OUT") + fmt.Sprintf("%s", PathSeparator)
	SslCrt = mapSslConf.String("CRT")
	SslKey = mapSslConf.String("KEY")
	return this
}

// 中间件装载器
func (this *Higo) Middleware(imiddleware ...IMiddleware) *Higo {
	for _, middleware := range imiddleware {
		this.middle = append(this.middle, middleware)
	}
	return this
}

// http服务
func (this *Higo) HttpServe(conf string, router IRouterLoader) *Higo {
	this.serve = append(this.serve, Hse{Config: conf, Router: router, Serve: "http"})
	return this
}

// https服务
func (this *Higo) HttpsServe(conf string, router IRouterLoader) *Higo {
	this.serve = append(this.serve, Hse{Config: conf, Router: router, Serve: "https"})
	return this
}

// websocket服务
func (this *Higo) WebsocketServe(conf string, router IRouterLoader) *Higo {
	this.serve = append(this.serve, Hse{Config: conf, Router: router, Serve: "websocket"})
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
	// 服务
	for _, s := range this.serve {
		// 初始化、加载配置、路由
		hg := Init().LoadConfigur(this.GetRoot())
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
		if gin.ReleaseMode == configure.ValueToString("MODE") {
			gin.SetMode(gin.ReleaseMode)
		}

		configs := configure.Config(s.Config)
		addr := configs.String("Addr")
		readTimeout := configs.Int("ReadTimeout")
		writeTimeout := configs.Int("WriteTimeout")

		handler := s.Router.Loader(hg)
		hg.loadRoute()

		serve := &http.Server{
			Addr:         configs.String("Addr"),
			Handler:      handler,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
		}

		router.Clear() //初始化路由容器

		if s.Serve == "http" {
			this.errgroup.Go(func() error {
				logger.Logrus.Infoln("HTTP Server listening at " + addr + " Starting Success!")
				return serve.ListenAndServe()
			})
		}
		if s.Serve == "https" {
			this.errgroup.Go(func() error {
				logger.Logrus.Infoln("HTTPS Server listening at " + addr + " Starting Success!")
				return serve.ListenAndServeTLS(SslOut+SslCrt, SslOut+SslKey)
			})
		}
		if s.Serve == "websocket" {
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
	router.AddRoute("", relativePath, "", router.IsStatic(true))
	hg.Engine.StaticFile(relativePath, filepath)
	return this
}

// 装载路由
func (this *Higo) loadRoute() *Higo {
	for _, route := range router.GetRoutes() {
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
	}
	return this
}

// 路由组Handle
func (this *Higo) GroupHandle(route *router.Route) *Higo {
	if handle := Convert(route.Handle()); handle != nil {
		handles := handleSlice(route)
		handles = append(handles, handle)
		this.group.Handle(strings.ToUpper(route.Method()), route.RelativePath(), handles...)
	}
	return this
}

// 路由Handle
func (this *Higo) Handle(route *router.Route) *Higo {
	if handle := Convert(route.Handle()); handle != nil {
		handles := handleSlice(route)
		handles = append(handles, handle)
		this.Engine.Handle(strings.ToUpper(route.Method()), route.RelativePath(), handles...)
	}
	return this
}

func handleSlice(route *router.Route) []gin.HandlerFunc {
	handles := make([]gin.HandlerFunc, 0)
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
		throw.Throw(throw.Message(err), throw.Code(0))
	}
	return this
}
