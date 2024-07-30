package higo

import (
	"github.com/dunpju/higo-config/config"
	"github.com/dunpju/higo-router/router"
	"github.com/dunpju/higo-utils/utils/sliceutil"
	"github.com/dunpju/higo-wsock/wsock"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"reflect"
	"sync"
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
	WsContainer      *wsock.WebsocketClient
	refWsResponder   reflect.Type
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
		WsContainer = wsock.NewWebsocketClient()
		refWsResponder = reflect.TypeOf((WebsocketResponder)(nil))
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
