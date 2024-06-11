package higo

import (
	"sync"
	"time"
)

var lock = &Locker{
	key: &sync.Map{},
}

type Locker struct {
	key *sync.Map
}

func (this *Locker) Lock(key string, fn func()) {
	_, ok := this.key.Load(key)
	if !ok {
		this.key.Store(key, key)
		defer this.UnLock(key)
		fn()
	}
}

func (this *Locker) Try(interval time.Duration, retry int) {

}

func (this *Locker) UnLock(key string) {
	this.key.Delete(key)
}

func Lock(key string, fn func()) {
	lock.Lock(key, fn)
}

func TryLock(key string) {

}

func UnLock(key string) {
	lock.UnLock(key)
}
