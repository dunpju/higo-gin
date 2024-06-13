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
	Key               string
	Timeout, Interval time.Duration
	Retry, counter    int
}

func (this *Mutex) Lock(locker *Locker, task func()) bool {
	_, ok := this.key.LoadOrStore(locker.Key, &sync.Mutex{})
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

func (this *Mutex) Retry(locker *Locker, task func()) bool {
retry:
	if !this.Lock(locker, task) {
		locker.counter++
		time.Sleep(locker.Interval)
		if locker.counter <= locker.Retry {
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
	if locker.Retry > 0 {
		return lock.Retry(locker, task)
	}
	return lock.Lock(locker, task)
}

func UnLock(key string) {
	lock.UnLock(key)
}
