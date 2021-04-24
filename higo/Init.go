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
)

var (
	initOnce             sync.Once
	serves               []*Serve
	onlySupportServe     *router.UniqueString
	pathSeparator        string
	root                 *utils.SliceString
	Upgrader             websocket.Upgrader
	WebsocketPingHandler WebsocketPingFunc
	WebsocketContainer   *WebsocketClient
	reflectWsResponder   reflect.Type
	WsCheckOrigin        = func(r *http.Request) bool {
		return true
	}
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
		Upgrader = websocket.Upgrader{
			CheckOrigin: WsCheckOrigin,
		}
		WebsocketPingHandler = websocketPingFunc
		WebsocketContainer = NewWebsocketClient()
		reflectWsResponder = reflect.TypeOf((WebsocketResponder)(nil))
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

func websocketPingFunc(websocketConn *WebsocketConn, waittime time.Duration) {
	time.Sleep(waittime)
	err := websocketConn.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
	if err != nil {
		WebsocketContainer.Remove(websocketConn.conn)
		return
	}
}
