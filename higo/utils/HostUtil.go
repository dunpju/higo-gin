package utils

import (
	"higo.yumi.com/src/app/Consts"
	"regexp"
	"strings"
)

// 获取http地址
func HttpAddr(requestHost string) string  {

	// 正则判断是否是ip:port
	if m, _ := regexp.MatchString(Consts.IP_PORT, requestHost); !m {
		return requestHost
	}

	// 分割
	requestHostSplit := strings.Split(requestHost, ":")

	return requestHostSplit[0]
}
