package higo

import (
	"github.com/dengpju/higo-config/config"
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-utils/utils/sliceutil"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
	"net/http"
	"reflect"
	"sync"
	"time"
)

const (
	ROOT           = "./../"
	HttpServe      = "http"
	HttpsServe     = "https"
	WebsocketServe = "websocket"
	WsConnIp       = "ws_conn_ip"
	WsRespstring   = "string"
	WsRespmap      = "map"
	WsRespstruct   = "struct"
	WsResperror    = "error"
	WsRespclose    = "close"
)

var (
	initOnce         sync.Once
	serves           []*Serve
	onlySupportServe *router.UniqueString
	pathSeparator    string
	AppConfigDir     *sliceutil.SliceString
	root             *sliceutil.SliceString
	MiddleCorsFunc   func(hg *Higo) gin.HandlerFunc
	MiddleAuthFunc   func(hg *Higo) gin.HandlerFunc
	Upgrader         websocket.Upgrader
	WsPingHandle     WebsocketPingFunc
	WsContainer      *WebsocketClient
	refWsResponder   reflect.Type
	WsCheckOrigin    WebsocketCheckFunc
	WsPitpatSleep    time.Duration
)

func init() {
	initOnce.Do(func() {
		serves = make([]*Serve, 0)
		container = NewDependency()
		RouterContainer = NewRouterCollect()
		taskList = make(chan *TaskExecutor)
		taskCron = cron.New(cron.WithSeconds())
		onlySupportServe = router.NewUniqueString()
		onlySupportServe.
			Append(HttpServe).
			Append(HttpsServe).
			Append(WebsocketServe)
		root = sliceutil.NewSliceString(".", "..", "")
		AppConfigDir = sliceutil.NewSliceString()
		MiddleCorsFunc = middleCorsFunc
		MiddleAuthFunc = middleAuthFunc
		WsCheckOrigin = func(r *http.Request) bool {
			return true
		}
		Upgrader = websocket.Upgrader{
			CheckOrigin: WsCheckOrigin,
		}
		WsPingHandle = wsPingFunc
		WsContainer = NewWebsocketClient()
		refWsResponder = reflect.TypeOf((WebsocketResponder)(nil))
		WsPitpatSleep = time.Second * 1
		config.AppPrefix = "config"
		config.AuthPrefix = config.AppPrefix
		config.AnnoPrefix = config.AppPrefix
		config.DbPrefix = config.EnvConf
		config.ServePrefix = config.EnvConf
	})

	chlist := getTaskList()
	go func() {
		for t := range chlist {
			doTask(t)
		}
	}()
}

func Root() *sliceutil.SliceString {
	return root
}
