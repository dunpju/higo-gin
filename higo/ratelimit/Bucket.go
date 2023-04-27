package ratelimit

import (
	"sync"
	"time"
)

type Bucket struct {
	cap      int64 //令牌数量
	tokens   int64 // token数量
	lock     sync.Mutex
	rate     int64 // 令牌速率，每秒生成多少个令牌
	lastTime int64
}

func NewBucket(cap, rate int64) *Bucket {
	if cap <= 0 || rate <= 0 {
		panic("error cap")
	}
	bucket := &Bucket{cap: cap, tokens: cap, rate: rate}
	//bucket.start() // 使用了更优雅的方式生成token
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
	now := time.Now().Unix()
	// 计算是否生成token
	this.tokens = this.tokens + (now-this.lastTime)*this.rate
	if this.tokens > this.cap {
		this.tokens = this.cap
	}
	this.lastTime = now
	if this.tokens > 0 {
		this.tokens--
		return true
	}
	return false
}
