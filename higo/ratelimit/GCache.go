package ratelimit

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type cacheData struct {
	key    string
	value  interface{}
	expire time.Time
}

func newCacheData(key string, value interface{}, expire time.Time) *cacheData {
	return &cacheData{key: key, value: value, expire: expire}
}

type GCache struct {
	max   int
	elist *list.List
	edata map[string]*list.Element
	lock  sync.Mutex
}

type GCacheOption func(g *GCache)
type GCacheOptions []GCacheOption

func (this GCacheOptions) apply(g *GCache) {
	for _, fn := range this {
		fn(g)
	}
}

func WithMax(size int) GCacheOption {
	return func(g *GCache) {
		if size > 0 {
			g.max = size
		}
	}
}

func NewGCache(opt ...GCacheOption) *GCache {
	cache := &GCache{elist: list.New(), edata: make(map[string]*list.Element, 0), max: 0}
	GCacheOptions(opt).apply(cache)
	cache.clear()
	return cache
}

func (this *GCache) clear() {
	go func() {
		for {
			this.removeExpired()
			time.Sleep(time.Second * 1)
		}
	}()
}

func (this *GCache) Get(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	if v, ok := this.edata[key]; ok {
		if time.Now().After(v.Value.(*cacheData).expire) {
			this.removeItem(v)
			return nil
		}
		this.elist.MoveToFront(v)
		return v.Value.(*cacheData).value
	}
	return nil
}

const NotExpireTTL = time.Hour * 24 * 365

func (this *GCache) Set(key string, newv interface{}, ttl time.Duration) {
	this.lock.Lock()
	defer this.lock.Unlock()
	var setExpire time.Time
	if ttl == 0 {
		setExpire = time.Now().Add(NotExpireTTL)
	} else {
		setExpire = time.Now().Add(ttl)
	}
	newCache := newCacheData(key, newv, setExpire)
	if v, ok := this.edata[key]; ok {
		v.Value = newCache
		this.elist.MoveToFront(v)
	} else {
		this.edata[key] = this.elist.PushFront(newCache)
		if this.max > 0 && len(this.edata) > this.max {
			this.removeOldest()
		}
	}
}

func (this *GCache) Print() {
	ele := this.elist.Front()
	if ele == nil {
		return
	}
	for {
		fmt.Println(this.Get(ele.Value.(*cacheData).key))
		ele = ele.Next()
		if ele == nil {
			break
		}
	}
}

func (this *GCache) removeOldest() {
	back := this.elist.Back()
	if back == nil {
		return
	}
	this.removeItem(back)
}

func (this *GCache) removeItem(els *list.Element) {
	key := els.Value.(*cacheData).key
	delete(this.edata, key)
	this.elist.Remove(els)
}

func (this *GCache) removeExpired() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, v := range this.edata {
		if time.Now().After(v.Value.(*cacheData).expire) {
			this.removeExpired()
		}
	}
}

func (this *GCache) Len() int {
	return len(this.edata)
}
