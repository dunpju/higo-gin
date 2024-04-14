package Codes

import (
	"gitee.com/dengpju/higo-code/code"
	. "github.com/dunpju/higo-gin/higo/errcode"
)

const (
	NotFound ErrorCode = iota + 400  //没找到
)

func code400() {
	code.Container().
	    Put(NotFound, "没找到")
}