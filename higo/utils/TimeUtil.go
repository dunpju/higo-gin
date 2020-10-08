package utils

import "time"

// 获取当前时间戳
func CurrentTimestamp () int64  {
	timestamp := time.Now().Unix()
	return int64(int(timestamp))
}

// 时间戳转时间
func TimestampToTime(ts int64) string {
	return time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05")
}

// 纳秒
func Nanoseconds(LowDateTime uint32, HighDateTime uint32) int64 {
	// 100-nanosecond intervals since January 1, 1601
	nsec := int64(HighDateTime)<<32 + int64(LowDateTime)
	// change starting time to the Epoch (00:00:00 UTC, January 1, 1970)
	nsec -= 116444736000000000
	// convert into nanoseconds
	nsec *= 100
	return nsec
}