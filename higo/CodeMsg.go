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

var codeMsg = map[int]*CodeMsg{}

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
	var cm *CodeMsg
	if len(codeMsg) > 0 {
		// 存在容器中
		cm, ok := codeMsg[code]
		if ok {
			return cm
		}
	}
	cm = &CodeMsg{code, msg}
	codeMsg[code] = cm
	return cm
}
