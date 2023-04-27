package ratelimit

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type LimiterCache struct {
	data sync.Map
}

var (
	IpCache *GCache
)

func init() {
	IpCache = NewGCache(WithMax(10000))
}

func IpLimiter(cap, rate int64, key string) func(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			ip := ClientIP(ctx.Request)
			var limiter *Bucket
			if v := IpCache.Get(ip); v != nil {
				limiter = v.(*Bucket)
			} else {
				limiter = NewBucket(cap, rate)
				IpCache.Set(ip, limiter, time.Second*5)
			}
			if ctx.Query(key) != "" {
				if limiter.IsAccept() {
					handler(ctx)
				} else {
					ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "too many requests-param"})
				}
			} else {
				handler(ctx)
			}
		}
	}
}

func ClientIP(r *http.Request) string {
	ip := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
