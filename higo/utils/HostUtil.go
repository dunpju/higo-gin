package utils

import (
	"github.com/dengpju/higo-gin/higo/consts"
	"regexp"
	"strings"
)

// 获取http地址
func HttpAddr(requestHost string) string  {

	// 正则判断是否是ip:port
	if m, _ := regexp.MatchString(consts.IP_PORT, requestHost); !m {
		return requestHost
	}

	// 分割
	requestHostSplit := strings.Split(requestHost, ":")

	return requestHostSplit[0]
}
