package higo

type TaskFunc func(params ...interface{})

var taskList chan *TaskExecutor

func getTaskList() chan *TaskExecutor{
	return taskList
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