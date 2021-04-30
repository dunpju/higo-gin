package higo

import (
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-utils/utils"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
	"net/http"
	"reflect"
	"sync"
	"time"
)

const (
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
	root             *utils.SliceString
	Upgrader         websocket.Upgrader
	WsPingHandle     WebsocketPingFunc
	WsContainer      *WebsocketClient
	refWsResponder   reflect.Type
	WsCheckOrigin    WebsocketCheckFunc
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
		WsCheckOrigin = func(r *http.Request) bool {
			return true
		}
		Upgrader = websocket.Upgrader{
			CheckOrigin: WsCheckOrigin,
		}
		WsPingHandle = wsPingFunc
		WsContainer = NewWebsocketClient()
		refWsResponder = reflect.TypeOf((WebsocketResponder)(nil))
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
