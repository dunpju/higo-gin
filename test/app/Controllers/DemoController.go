package Controllers

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
)

// 测试异常
func HttpsTestThrow(ctx *gin.Context) string  {
	higo.Throw("https 测试异常", 0)
	return "https_test_throw"
}

// 测试get请求
func HttpsTestGet(ctx *gin.Context) string  {
	return "https_test_get"
}

// 测试post请求
func HttpsTestPost(ctx *gin.Context) string {
	return "https_test_post"
}

// 测试异常
func HttpTestThrow(ctx *gin.Context) string  {
	higo.Throw("http 测试异常", 0)
	return "http_test_throw"
}

// 测试get请求
func HttpTestGet(ctx *gin.Context) string  {
	return "http_test_get"
}

// 测试post请求
func HttpTestPost(ctx *gin.Context) string {
	return "http_test_post"
}