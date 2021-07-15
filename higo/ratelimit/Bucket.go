package ratelimit

import (
	"sync"
	"time"
)

type Bucket struct {
	cap    int //令牌数量
	tokens int // token数量
	lock   sync.Mutex
	rate   int // 令牌速率，每秒生成多少个令牌
}

func NewBucket(cap, rate int) *Bucket {
	if cap <= 0 || rate <= 0 {
		panic("error cap")
	}
	bucket := &Bucket{cap: cap, tokens: cap, rate: rate}
	bucket.start()
	return bucket
}

func (this *Bucket) start() {
	go func() {
		for {
			time.Sleep(time.Second)
			this.addToken()
		}
	}()
}

func (this *Bucket) addToken() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.tokens+this.rate <= this.cap {
		this.tokens += this.rate
	} else {
		this.tokens = this.cap
	}
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
