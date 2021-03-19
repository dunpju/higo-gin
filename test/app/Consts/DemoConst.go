package Consts

// CODE 常量,格式：code_TO_msg
const (
	// `vv`
	SUCCESS          = "20000@成功"
	SERVER_ERROR     = "50000@系统错误"
	AUTH_ERROR       = "50001@认证错误"
	UNKNOWN_ERROR    = "50002@未知错误"
	RSA_ERROR        = "50003@解密错误"
	PARAMETER_ERROR  = "50004@参数错误"
	NOT_FOUND        = "50005@未找到"
	CODE_ERROR       = "50006@失败"
	INVALID_TOKEN    = "50007@无效token"
	INVALID_API      = "50008@无效api"
	INVALID_MAP      = "50009@无效api映射"
	// 业务异常从51000开始
	TEST_ERROR       = "51000@测试异常"
)
