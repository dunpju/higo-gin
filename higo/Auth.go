package higo

// 是否空标记
func IsEmptyFlag(route Route)  {
	if route.Flag == "" {
		Throw(route.RelativePath + "未设置标记",0)
	}
}

// 是否不用鉴权
func IsNotAuth(flag string) bool {
	if "" == flag {
		return false
	}
	// 空配置
	if nil == Container().Configure {
		return false
	}
	notAuth := Container().Config("NotAuth")
	_, ok := notAuth[flag]
	return ok
}
