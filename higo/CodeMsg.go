package higo

import (
	"strconv"
	"strings"
)

// code 对应 msg
type CodeMsg struct {
	Code int
	Msg  string
}

var container = map[int]*CodeMsg{}

// 常量函数
func Const(constCode string) *CodeMsg {
	return separate(constCode)
}

// _TO_分割
func separate(constCode string) *CodeMsg {
	constSlice := strings.Split(constCode, "_TO_")
	code, err := strconv.Atoi(constSlice[0])
	msg := constSlice[1]
	if err != nil {
		codeError := strings.Split("0_TO_常量格式错误", "_TO_")
		code, _ = strconv.Atoi(constSlice[0])
		msg = codeError[1]
	}
	var codeMsg *CodeMsg
	if len(container) > 0 {
		// 存在容器中
		codeMsg, ok := container[code]
		if ok {
			return codeMsg
		}
	}
	codeMsg = &CodeMsg{code, msg}
	container[code] = codeMsg
	return codeMsg
}
