package higo

func init() {
	Once.Do(func() {
		config = make(Configure)
		container = make(Dependency)
		Router = make(RouterCollect)
	})
}
