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
	return !ok
}

func (this *Mutex) Lock(key string, fn func()) bool {
	return this.lock(key, fn)
}

func (this *Mutex) Retry(returner *Returner, key string, fn func()) bool {
start:
	if !this.lock(key, fn) {
		returner.counter++
		time.Sleep(returner.Interval)
		if returner.counter <= returner.Retry {
			goto start
		}
		return false
	}
	return true
}

func (this *Mutex) UnLock(key string) {
	this.key.Delete(key)
}

func Lock(key string, fn func()) bool {
	return lock.Lock(key, fn)
}

type Returner struct {
	Interval       time.Duration
	Retry, counter int
}

func Retry(returner *Returner, key string, fn func()) bool {
	return lock.Retry(returner, key, fn)
}

func UnLock(key string) {
	lock.UnLock(key)
}
