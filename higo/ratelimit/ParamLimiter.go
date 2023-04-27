package ratelimit

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ParamLimiter(cap, rate int64, key string) func(handler gin.HandlerFunc) gin.HandlerFunc {
	limiter := NewBucket(cap, rate)
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(ctx *gin.Context) {
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
