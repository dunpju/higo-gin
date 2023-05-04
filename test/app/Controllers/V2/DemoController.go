package V2

import (
	"github.com/dunpju/higo-throw/exception"
	"github.com/gin-gonic/gin"
)

// 测试异常
func HttpsTestThrow(ctx *gin.Context) string {
	exception.Throw(exception.Message("v2 https 测试异常"), exception.Code(0))
	return "v2 https_test_throw"
}

// 测试get请求
func HttpsTestGet(ctx *gin.Context) string {
	return "v2 https_test_get"
}

// 测试post请求
func HttpsTestPost(ctx *gin.Context) string {
	return "v2 https_test_post"
}

// 测试异常
func HttpTestThrow(ctx *gin.Context) string {
	exception.Throw(exception.Message("v2 http 测试异常"), exception.Code(0))
	return "v2 http_test_throw"
}

// 测试get请求
func HttpTestGet(ctx *gin.Context) string {
	return "v2 http_test_get"
}

// 测试post请求
func HttpTestPost(ctx *gin.Context) string {
	return "v2 http_test_post"
}
