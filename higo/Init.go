package higo

import (
	"github.com/dengpju/higo-router/router"
	"github.com/robfig/cron/v3"
)

const (
	HttpServe      = "http"
	HttpsServe     = "https"
	WebsocketServe = "websocket"
)

var (
	serves           []*Serve
	onlySupportServe *router.UniqueString
)

func init() {
	Once.Do(func() {
		serves = make([]*Serve, 0)
		container = make(Dependency)
		RouterContainer = make(RouterCollect)
		taskList = make(chan *TaskExecutor)
		taskCron = cron.New(cron.WithSeconds())
		onlySupportServe = router.NewUniqueString()
		onlySupportServe.Append(HttpServe).
			Append(HttpsServe).
			Append(WebsocketServe)
	})

	chlist := getTaskList()
	go func() {
		for t := range chlist {
			doTask(t)
		}
	}()
}
