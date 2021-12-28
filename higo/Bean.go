package higo

type Bean struct {
	Middleware
}

func NewBean() *Bean {
	return &Bean{}
}

func (this *Bean) Provider() {}

func (this *Bean) NewServe(conf string) *Serve {
	return NewServe(conf)
}

// 使用权交由用户决定
//func (this *Bean) NewOrm() *Orm {
//	return NewOrm()
//}

// 使用权交由用户决定
//func (this *Bean) NewRedisPool() *redis.Pool {
//	return RedisPool
//}

// 使用权交由用户决定
//func (this *Bean) NewRedisAdapter() *RedisAdapter {
//	return NewRedisAdapter()
//}
