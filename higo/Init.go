package higo

import (
	"github.com/robfig/cron/v3"
)

func init() {
	Once.Do(func() {
		container = make(Dependency)
		RouterContainer = make(RouterCollect)
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
