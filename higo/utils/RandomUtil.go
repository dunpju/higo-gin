package utils

import (
	"higo.yumi.com/src/app/Consts"
	"math/rand"
	"time"
)

type RandomUtil struct {}

func NewRandomUtil() *RandomUtil {
	return &RandomUtil{}
}

// 随机0-86400
func (this *RandomUtil) IntHour24ToSecond() int64 {
	return rand.Int63n(Consts.HOUR_24_TO_SECOND)
}

// 随机字符串
func (this *RandomUtil) Str(length int) string  {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
