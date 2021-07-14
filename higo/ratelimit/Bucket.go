package ratelimit

import "sync"

type Bucket struct {
	cap    int //令牌数量
	tokens int // token数量
	lock   sync.Mutex
}

func NewBucket(cap int) *Bucket {
	if cap <= 0 {
		panic("error cap")
	}
	return &Bucket{cap: cap, tokens: cap}
}

func (this *Bucket) IsAccept() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.tokens > 0 {
		this.tokens--
		return true
	}
	return false
}
