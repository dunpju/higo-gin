package higo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
	"higo.yumi.com/src/higo/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var (
	hg *Higo
	// ssl 证书
	SslOut, SslCrt, SslKey string
)

// 上下文
type IHiContext interface {
	OnRequest(*gin.Context) error
	OnResponse(result interface{}) (interface{}, error)
}

type Higo struct {
	*gin.Engine
	g            *gin.RouterGroup
	eg           errgroup.Group
	exprData     map[string]interface{}
	currentGroup string
	root         string
	containers   *Containers
}

// 初始化
func Init() *Higo {
	hg = &Higo{
		Engine:     gin.New(),
		exprData:   map[string]interface{}{},
		containers: NewContainer(),
	}
	// 全局异常
	hg.Engine.Use(NewRecover().RuntimeException(hg))

	// 获取主目录
	root := hg.GetRoot()

	// runtime目录
	runtime := root + "runtime"
	if _, err := os.Stat(runtime); os.IsNotExist(err) {
		if os.Mkdir(runtime, os.ModePerm) != nil {}
	}

	// 日志
	Log(root)

	// 装载配置
	filepathErr := filepath.Walk(root + "conf",
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			if path.Ext(p) == ".yaml" {
				fmt.Println("yaml file:", filepath.Base(p))
				yamlFile, _ := ioutil.ReadFile(p)
				yamlFileErr := yaml.Unmarshal(yamlFile, &Container().Configure)
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
	SslOut = root + mapSslConf["OUT"].(string)
	SslCrt = mapSslConf["CRT"].(string)
	SslKey = mapSslConf["KEY"].(string)
	// 生成ssl证书
	utils.NewSsl(SslOut, SslCrt, SslKey).Generate()

	return hg
}

// 设置主目录
func (this *Higo) SetRoot(root string) {
	this.root = root
}

// 获取主目录
func (this *Higo) GetRoot() string {
	return utils.If(this.root == "", ROOT, this.root).(string)
}

// 中间件装载器
func (this *Higo) Middleware(imiddleware ...IMiddleware) *Higo {
	for _, middleware := range imiddleware {
		this.Engine.Use(middleware.Loader(this))
	}
	return this
}

// http服务
func (this *Higo) HttpServe(conf string, router IRouterLoader) *Higo {
	config := Container().Config(conf)
	addr, _ := config["Addr"]
	rt, _ := config["ReadTimeout"]
	wt, _ := config["WriteTimeout"]
	readTimeout, _ := rt.(int)
	writeTimeout, _ := wt.(int)
	httpServe := &http.Server{
		Addr:         addr.(string),
		Handler:      router.Loader(this),
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}
	this.eg.Go(func() error {
		return httpServe.ListenAndServe()
	})
	return this
}

// https服务
func (this *Higo) HttpsServe(conf string, router IRouterLoader) *Higo {
	config := Container().Config(conf)
	addr, _ := config["Addr"]
	rt, _ := config["ReadTimeout"]
	wt, _ := config["WriteTimeout"]
	readTimeout, _ := rt.(int)
	writeTimeout, _ := wt.(int)
	httpsServe := &http.Server{
		Addr:         addr.(string),
		Handler:      router.Loader(this),
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}
	this.eg.Go(func() error {
		return httpsServe.ListenAndServeTLS(SslOut + SslCrt, SslOut + SslKey)
	})
	return this
}

// websocket服务
func (this *Higo) WebsocketServe(conf string, router IRouterLoader) *Higo {
	config := Container().Config(conf)
	addr, _ := config["Addr"]
	rt, _ := config["ReadTimeout"]
	wt, _ := config["WriteTimeout"]
	readTimeout, _ := rt.(int)
	writeTimeout, _ := wt.(int)
	websocket := &http.Server{
		Addr:         addr.(string),
		Handler:      router.Loader(this),
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}
	this.eg.Go(func() error {
		return websocket.ListenAndServe()
	})
	return this
}

// 启动
func (this *Higo) Boot() {
	fmt.Println("启动成功")
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
	return Container().GetRoute(relativePath), true
}

// 路由组
func (this *Higo) AddGroup(prefix string, routes ...Route) *Higo {
	this.g = this.Engine.Group(prefix)
	for _, route := range routes {
		// 判断空标记
		IsEmptyFlag(route)
		// 添加路由容器
		Container().AddRoutes(route.RelativePath, route)
		method := strings.ToUpper(route.Method)
		this.GroupHandle(method, route.RelativePath, route.Handle)
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
	fmt.Printf("%T\n",handler)
	if h := Convert(handler); h != nil {
		this.Engine.Handle(httpMethod, relativePath, h)
	}
	return this
}

func (this *Higo) Mount(group string, icontroller ...IController) *Higo {
	this.g = this.Engine.Group(group)
	for _, controller := range icontroller {
		this.currentGroup = group
		controller.Controller(this)
	}
	return this
}