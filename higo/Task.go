package higo

import "github.com/robfig/cron/v3"

type TaskFunc func(params ...interface{})

var taskList chan *TaskExecutor
var taskCron *cron.Cron

func getTaskList() chan *TaskExecutor{
	return taskList
}

func getCronTask() *cron.Cron {
	return taskCron
}

type TaskExecutor struct {
	fn TaskFunc
	params []interface{}
	callback func()
}

func NewTaskExecutor(fn TaskFunc, params []interface{}, callback func()) *TaskExecutor {
	return &TaskExecutor{fn: fn, params: params, callback: callback}
}

func (this *TaskExecutor) Exec()  {
	this.fn(this.params...)
}

func Task(fn TaskFunc, callback func(), params ...interface{})  {
	if fn == nil {
		return
	}
	go func() {
		getTaskList() <- NewTaskExecutor(fn, params, callback)
	}()
}

func doTask(t * TaskExecutor)  {
	go func() {
		defer func() {
			if t.callback != nil {
				t.callback()
			}
		}()
		t.Exec()
	}()
}