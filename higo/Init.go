package higo

import (
	"github.com/dengpju/higo-router/router"
	"github.com/dengpju/higo-utils/utils"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
	"net/http"
	"sync"
)

const (
	HttpServe      = "http"
	HttpsServe     = "https"
	WebsocketServe = "websocket"
)

var (
	initOnce                 sync.Once
	serves                   []*Serve
	onlySupportServe         *router.UniqueString
	pathSeparator            string
	root                     *utils.SliceString
	Upgrader                 websocket.Upgrader
	WebsocketPongHandler     WebsocketPongFunc
	WebsocketClientContainer *WebsocketClient
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
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		WebsocketPongHandler = websocketPongFunc
		WebsocketClientContainer = NewWebsocketClient()
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
