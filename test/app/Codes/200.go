package Codes

import (
	"gitee.com/dengpju/higo-code/code"
	. "github.com/dunpju/higo-gin/higo/errcode"
)

const (
	Success ErrorCode = iota + 200 //成功
)

func code200() {
	code.Container().
		Put(Success, "成功")
}
