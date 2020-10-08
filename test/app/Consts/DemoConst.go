package Consts

// CODE 常量,格式：code_TO_msg
const (
	SUCCESS          = "20000_TO_成功"
	SERVER_ERROR     = "50000_TO_系统错误"
	AUTH_ERROR       = "50001_TO_认证错误"
	UNKNOWN_ERROR    = "50002_TO_未知错误"
	RSA_ERROR        = "50003_TO_解密错误"
	PARAMETER_ERROR  = "50004_TO_参数错误"
	NOT_FOUND        = "50005_TO_未找到"
	CODE_ERROR       = "50006_TO_失败"
	INVALID_TOKEN    = "50007_TO_无效token"
	INVALID_API      = "50008_TO_无效api"
	INVALID_MAP      = "50009_TO_无效api映射"
	// 业务异常从51000开始
	TEST_ERROR       = "51000_TO_测试异常"
)
