package higo

import (
	"sync"
	"time"
)

var lock = &Mutex{
	key: &sync.Map{},
}

type Mutex struct {
	key *sync.Map
}

func (this *Mutex) lock(key string, fn func()) bool {
	mutex, ok := this.key.LoadOrStore(key, &sync.Mutex{})
	if !ok {
		mutex.(*sync.Mutex).Lock()
		defer mutex.(*sync.Mutex).Unlock()
		defer this.UnLock(key)
		fn()
	}
	return ok
}

func (this *Mutex) Lock(key string, fn func()) bool {
	return !this.lock(key, fn)
}

func (this *Mutex) Retry(interval time.Duration, retry int, key string, fn func()) {
	counter := 0
start:
	if this.lock(key, fn) {
		counter++
		time.Sleep(interval)
		if counter <= retry {
			goto start
		}
	}
}

func (this *Mutex) UnLock(key string) {
	this.key.Delete(key)
}

func Lock(key string, fn func()) {
	lock.Lock(key, fn)
}

func Retry(interval time.Duration, retry int, key string, fn func()) {
	lock.Retry(interval, retry, key, fn)
}

func UnLock(key string) {
	lock.UnLock(key)
}
