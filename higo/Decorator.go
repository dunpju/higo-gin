package higo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/ffjson/ffjson"
)

func CacheDecorator(h gin.HandlerFunc, param string, format string, empty interface{}) gin.HandlerFunc {
	return func(context *gin.Context) {
		getId := context.Param(param)
		key := fmt.Sprintf(format, getId)
		ret, err := Redis.GetByte(key)
		if err != nil {
			h(context)
			dbResult, exists := context.Get("db_result")
			if !exists {
				dbResult = empty
			}
			retData, _ := ffjson.Marshal(dbResult)
			Redis.Setex(key, retData, 20)
			context.JSON(200, retData)
		} else {
			_ = ffjson.Unmarshal(ret, &empty)
			context.JSON(200, empty)
		}
	}
}
