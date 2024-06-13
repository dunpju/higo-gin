package higo

import (
	"context"
	"sync"
	"time"
)

var lock = &Mutex{
	key: &sync.Map{},
}

type Mutex struct {
	key *sync.Map
}

type Locker struct {
	Key     string
	Timeout time.Duration
}

func (this *Mutex) Lock(locker *Locker, task func()) bool {
	_, ok := this.key.LoadOrStore(locker.Key, locker)
	if !ok {
		defer this.UnLock(locker.Key)
		if locker.Timeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), locker.Timeout)
			defer cancel()
			go func(ctx context.Context) {
				select {
				case <-ctx.Done():
					this.UnLock(locker.Key)
					return
				case <-time.After(locker.Timeout):
					this.UnLock(locker.Key)
					return
				}
			}(ctx)
		}
		task()
	}
	return !ok
}

func (this *Mutex) Retry(returner *Returner, locker *Locker, task func()) bool {
retry:
	if !this.Lock(locker, task) {
		returner.counter++
		time.Sleep(returner.Interval)
		if returner.counter <= returner.Retry {
			goto retry
		}
		return false
	}
	return true
}

func (this *Mutex) UnLock(key string) {
	this.key.Delete(key)
}

func Lock(locker *Locker, task func()) bool {
	return lock.Lock(locker, task)
}

type Returner struct {
	Interval       time.Duration
	Retry, counter int
}

func Retry(returner *Returner, locker *Locker, task func()) bool {
	return lock.Retry(returner, locker, task)
}

func UnLock(key string) {
	lock.UnLock(key)
}
