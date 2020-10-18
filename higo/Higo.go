package higo

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo/consts"
	"github.com/dengpju/higo-gin/higo/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	hg *Higo
	// 系统类型 Windows Or Linux
	SysType string
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
	g          *gin.RouterGroup
	eg         errgroup.Group
	root       string
	containers *Containers
	isAutoSsl  bool
	middle     []IMiddleware
	serve      []Hse
	attribute  []interface{}
}

// 初始化
func Init() *Higo {
	hg = &Higo{
		Engine:     gin.New(),
		containers: NewContainer(),
		middle:     make([]IMiddleware, 0),
		serve:      make([]Hse, 0),
	}

	// 全局异常
	hg.Engine.Use(NewRecover().Exception(hg))
	// 系统类型
	SysType = runtime.GOOS
	// 初始分隔符
	if SysType == "windows" {
		PathSeparator = "\\"
	} else {
		PathSeparator = "/"
	}
	// 是否使用自带ssl测试https
	hg.isAutoSsl = false

	return hg
}

// 设置主目录
func (this *Higo) SetRoot(root string) *Higo {
	this.root = root
	return this
}

// 获取主目录
func (this *Higo) GetRoot() string {
	return utils.If(this.root == "", consts.ROOT, this.root).(string)
}

// 配置
func (this *Higo) config() *Higo {
	// 获取主目录
	root := hg.GetRoot()
	// runtime目录
	runtimeDir := root + "runtime"
	if _, err := os.Stat(runtimeDir); os.IsNotExist(err) {
		if os.Mkdir(runtimeDir, os.ModePerm) != nil {}
	}
	// 日志
	Log(root)
	// 装载配置
	confDir := root + "conf"
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		if os.Mkdir(confDir, os.ModePerm) != nil {}
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
				fmt.Println("Loader configure file:", filepath.Base(p))
				yamlFile, _ := ioutil.ReadFile(p)
				yamlFileErr := yaml.Unmarshal(yamlFile, &Container().C)
				if yamlFileErr != nil {
					Throw(yamlFileErr,0)
				}
			}
			return nil
		})
	if filepathErr != nil {
		Throw(filepathErr,0)
	}
	mapSslConf := Container().Config("SSL")
	SslOut = root + mapSslConf["OUT"].(string) + fmt.Sprintf("%s", PathSeparator)
	SslCrt = mapSslConf["CRT"].(string)
	SslKey = mapSslConf["KEY"].(string)
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
func (this *Higo) IsAutoGenerateSsl(isAuto bool) *Higo {
	this.isAutoSsl = isAuto
	return this
}

// 启动
func (this *Higo) Boot() {
	// 服务
	for _, s := range this.serve {
		// 设置服务根目录
		hg := Init().SetRoot(this.GetRoot())
		// 配置
		hg.config()
		// 中间件
		for _, m := range this.middle {
			hg.Engine.Use(m.Loader(hg))
		}
		// 是否使用自带ssl测试https
		if this.isAutoSsl {
			// 生成ssl证书
			utils.NewSsl(SslOut, SslCrt, SslKey).Generate()
		}
		config := Container().Config(s.Config)
		addr, _ := config["Addr"]
		rt, _ := config["ReadTimeout"]
		wt, _ := config["WriteTimeout"]
		readTimeout, _ := rt.(int)
		writeTimeout, _ := wt.(int)
		serve := &http.Server{
			Addr:         addr.(string),
			Handler:      s.Router.Loader(hg),
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
		}

		if s.Serve == "http" {
			this.eg.Go(func() error {
				fmt.Println("HTTP Server listening at " + addr.(string) + " 启动成功\n")
				return serve.ListenAndServe()
			})
		}
		if s.Serve == "https" {
			this.eg.Go(func() error {
				fmt.Println("HTTPS Server listening at " + addr.(string) + " 启动成功\n")
				return serve.ListenAndServeTLS(SslOut + SslCrt, SslOut + SslKey)
			})
		}
		if s.Serve == "websocket" {
			this.eg.Go(func() error {
				fmt.Println("WEBSOCKET Server listening at " + addr.(string) + " 启动成功\n")
				return serve.ListenAndServe()
			})
		}
	}

	if err := this.eg.Wait(); err != nil {
		Logrus.Fatal(err)
	}
}

// 容器
func Container() *Containers {
	return hg.containers
}

// 获取路由
func (this *Higo) GetRoute(relativePath string) (Route, bool) {
	return Container().Route(relativePath), true
}

// 静态文件
func (this *Higo) StaticFile(relativePath, filepath string) *Higo {
	// 添加路由容器
	Container().AddRoutes(relativePath, Route{IsStatic: true})
	hg.Engine.StaticFile(relativePath, filepath)
	return this
}

// 路由组
func (this *Higo) AddGroup(prefix string, routes ...Route) *Higo {
	this.g = this.Engine.Group(prefix)
	for _, route := range routes {
		fmt.Printf("%T\n", route.Handle);
		// 判断空标记
		IsEmptyFlag(route)
		// 添加路由容器
		Container().AddRoutes("/" + strings.TrimLeft(prefix, "/") + "/" + strings.TrimLeft(route.RelativePath, "/"), route)
		method := strings.ToUpper(route.Method)
		this.GroupHandle(method, route.RelativePath, route.Handle)
		//RegisterDependencies(this)
	}
	return this
}

// 路由
func (this *Higo) AddRoute(routes ...Route) *Higo {
	for _, route := range routes {
		// 判断空标记
		IsEmptyFlag(route)
		// 添加路由容器
		Container().AddRoutes(route.RelativePath, route)
		method := strings.ToUpper(route.Method)
		this.Handle(method, route.RelativePath, route.Handle)
	}
	return this
}

// 路由组Handle
func (this *Higo) GroupHandle(httpMethod, relativePath string, handler interface{}) *Higo {
	if h := Convert(handler); h != nil {
		this.g.Handle(httpMethod, relativePath, h)
	}
	return this
}

// 路由Handle
func (this *Higo) Handle(httpMethod, relativePath string, handler interface{}) *Higo {
	if h := Convert(handler); h != nil {
		this.Engine.Handle(httpMethod, relativePath, h)
	}
	return this
}

/**
func (this *Higo) Mount(group string, icontroller ...IController) *Higo {
	this.g = this.Engine.Group(group)
	for _, controller := range icontroller {
		controller.Controller()
	}
	return this
}*/

//获取属性
func (this *Higo) getAttribute(t reflect.Type) interface{} {
	for _, p := range this.attribute {
		fmt.Println(t)
		fmt.Println(reflect.TypeOf(p))
		if t == reflect.TypeOf(p) {
			return p
		}
	}
	return nil
}


// 设置属性
func (this *Higo) setAttribute(builder IBuilder) {
	vClass := reflect.ValueOf(builder)
	vClassT := reflect.TypeOf(builder)
	if vClass.Kind() == reflect.Ptr {
		vClass = vClass.Elem()
	}
	vt := reflect.TypeOf(&Value{})
	fmt.Println(vt)
	for i := 0; i < vClass.NumField(); i++ {
		f := vClass.Field(i)
		fmt.Println(f)
		fmt.Println(f.Type())
		if vt != f.Type() {
			continue
		}
		if !f.IsNil() || f.Kind() != reflect.Ptr {
			continue
		}
		if p := this.getAttribute(f.Type()); p != nil {
			fmt.Println(111)
			f.Set(reflect.New(f.Type().Elem()))
			f.Elem().Set(reflect.ValueOf(p).Elem())
			if IsAnnotation(f.Type()) {
				p.(Annotation).SetTag(vClassT.Field(i).Tag)
			}
		}
	}
}

// 注册依赖
func (this *Higo) RegisterDependencies(builders ...IBuilder) *Higo {
	for _, builder := range builders {
		name := reflect.ValueOf(builder).Type().Name()
		Container().Di[name] = builder
		this.setAttribute(builder)
	}
	return this
}