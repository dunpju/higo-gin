package higo

import (
	"gitee.com/dengpju/higo-code/code"
	"github.com/dengpju/higo-config/config"
	"github.com/dengpju/higo-gin/test/app/Consts"
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-throw/exception"
	"github.com/dengpju/higo-utils/utils"
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
	AppConfigDir     *utils.SliceString
	root             *utils.SliceString
	MiddleCorsFunc   func(cxt *gin.Context)
	MiddleAuthFunc   func(cxt *gin.Context)
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
		container = make(Dependency)
		RouterContainer = make(RouterCollect)
		taskList = make(chan *TaskExecutor)
		taskCron = cron.New(cron.WithSeconds())
		onlySupportServe = router.NewUniqueString()
		onlySupportServe.
			Append(HttpServe).
			Append(HttpsServe).
			Append(WebsocketServe)
		root = utils.NewSliceString(".", "..", "")
		AppConfigDir = utils.NewSliceString()
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

func Root() *utils.SliceString {
	return root
}

func wsPingFunc(websocketConn *WebsocketConn, waittime time.Duration) {
	time.Sleep(waittime)
	err := websocketConn.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
	if err != nil {
		WsContainer.Remove(websocketConn.conn)
		return
	}
}

func middleCorsFunc(cxt *gin.Context) {
	method := cxt.Request.Method
	origin := cxt.Request.Header.Get("Origin") //请求头部
	if origin != "" {
		cxt.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		cxt.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		cxt.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		cxt.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		cxt.Header("Access-Control-Allow-Credentials", "true")
	}

	//允许类型校验
	if method == "OPTIONS" {
		cxt.AbortWithStatus(http.StatusNoContent)
	}
}

func middleAuthFunc(cxt *gin.Context) {
	if route, ok := hg.GetRoute(cxt.Request.URL.Path); ok {
		if ! IsNotAuth(route.Flag()) && !route.IsStatic() {
			if "" == cxt.GetHeader("X-Token") {
				exception.Throw(exception.Message(code.Message(Consts.InvalidToken).Message), exception.Code(code.Message(Consts.InvalidToken).Code))
			}
		}
	} else {
		exception.Throw(exception.Message(code.Message(Consts.InvalidApi).Message), exception.Code(code.Message(Consts.InvalidApi).Code))
	}
}
