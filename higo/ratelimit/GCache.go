package ratelimit

import (
	"container/list"
	"fmt"
	"sync"
)

type GCache struct {
	elist *list.List
	edata map[string]*list.Element
	lock  sync.Mutex
}

func NewGCache() *GCache {
	return &GCache{elist: list.New(), edata: make(map[string]*list.Element, 0)}
}

func (this *GCache) Get(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	if v, ok := this.edata[key]; ok {
		this.elist.MoveToFront(v)
		return v.Value
	}
	return nil
}

func (this *GCache) Set(key string, newv interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if v, ok := this.edata[key]; ok {
		v.Value = newv
		this.elist.MoveToFront(v)
	} else {
		this.edata[key] = this.elist.PushFront(newv)
	}
}

func (this *GCache) Print() {
	ele := this.elist.Front()
	if ele == nil {
		return
	}
	for {
		fmt.Println(ele.Value)
		ele = ele.Next()
		if ele == nil {
			break
		}
	}
}
