package higo

import (
	"strconv"
	"strings"
)

// code 对应 msg
type CodeMsg struct {
	Code int
	Msg string
}

// 获取函数
func Get(constCode string) *CodeMsg {
	return separate(constCode)
}

// 按照_TO_分割
func separate(constCode string) *CodeMsg {
	constSlice := strings.Split(constCode, "_TO_")
	code, err := strconv.Atoi(constSlice[0])
	msg := constSlice[1]
	if err != nil {
		codeError := strings.Split("0_TO_失败", "_TO_")
		code, _= strconv.Atoi(constSlice[0])
		msg = codeError[1]
	}
	return &CodeMsg{code, msg}
}
