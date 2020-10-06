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
