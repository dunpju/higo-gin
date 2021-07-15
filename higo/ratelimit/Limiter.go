package ratelimit

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Limiter(cap int) func(handler gin.HandlerFunc) gin.HandlerFunc {
	limiter := NewBucket(cap, 1)
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			if limiter.IsAccept() {
				handler(ctx)
			} else {
				ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "too many requests"})
			}
		}
	}
}
