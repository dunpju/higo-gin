package higo

func init() {
	Once.Do(func() {
		config = make(Configure)
		container = make(Dependency)
		Router = make(RouterCollect)
		taskList = make(chan *TaskExecutor)
	})

	chlist := getTaskList()
	go func() {
		for t := range chlist{
			doTask(t)
		}
	}()
}
