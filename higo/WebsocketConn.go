package higo

import (
	"gitee.com/dengpju/higo-code/code"
	"github.com/dunpju/higo-throw/exception"
	"github.com/dunpju/higo-utils/utils/maputil"
	"sync"
)

var (
	// WsRecoverHandle Recover处理函数(可自定义替换)
	WsRecoverHandle WsRecoverFunc
	wsRecoverOnce   sync.Once
)

func init() {
	wsRecoverOnce.Do(func() {
		WsRecoverHandle = func(r interface{}) (respMsg string) {
			if msg, ok := r.(*code.CodeMessage); ok {
				respMsg = maputil.Array().
					Put("code", msg.Code).
					Put("message", msg.Message).
					Put("data", nil).
					String()
			} else if arrayMap, ok := r.(maputil.ArrayMap); ok {
				respMsg = arrayMap.String()
			} else {
				respMsg = maputil.Array().
					Put("code", 0).
					Put("message", exception.ErrorToString(r)).
					Put("data", nil).
					String()
			}
			return
		}
	})
}

type WsRecoverFunc func(r interface{}) string
