package higo

type TaskFunc func(params ...interface{})

var taskList chan *TaskExecutor

func getTaskList() chan *TaskExecutor{
	return taskList
}


type TaskExecutor struct {
	fn TaskFunc
	params []interface{}
}

func (this *TaskExecutor) Exec()  {
	this.fn(this.params...)
}

func NewTaskExecutor(fn TaskFunc, params []interface{}) *TaskExecutor {
	return &TaskExecutor{fn: fn, params: params}
}

func Task(fn TaskFunc, params ...interface{})  {
	getTaskList() <- NewTaskExecutor(fn, params)
}