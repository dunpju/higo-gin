package higo

import (
	"github.com/robfig/cron/v3"
	"sync"
)

var redisOnce sync.Once

func init() {
	Once.Do(func() {
		config = make(Configure)
		container = make(Dependency)
		Router = make(RouterCollect)
		taskList = make(chan *TaskExecutor)
		taskCron = cron.New(cron.WithSeconds())
	})

	chlist := getTaskList()
	go func() {
		for t := range chlist{
			doTask(t)
		}
	}()
}
