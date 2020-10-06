package utils

import (
	"os"
	"syscall"
)

// 自定义文件结构体
type File struct {
	Name string // 文件名(完整路径)
}

// 构造函数
func NewFile(name string) *File {
	return &File{name}
}

// 获取文件创建时间戳
func (this *File) GetCreateTimestamp() int64 {
	fileInfo, _ := os.Stat(this.Name)

	stat_t := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	return stat_t.CreationTime.Nanoseconds()/1e9

	//sysType := runtime.GOOS
	//if sysType == "windows" {
	//	stat_t := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	//	return stat_t.CreationTime.Nanoseconds()/1e9
	//}
	//
	//// Sys()返回的是interface{}，所以需要类型断言，不同平台需要的类型不一样，linux上为*syscall.Stat_t;
	//// windows 上为 *syscall.Win32FileAttributeData
	//stat_t := fileInfo.Sys().(*syscall.Stat_t)
	//return stat_t.Ctim.Sec
}

// 获取文件创建时间
//func (this *File) GetCreateTime() string {
//	fileInfo, _ := os.Stat(this.Name)
//
//	sysType := runtime.GOOS
//	if sysType == "windows" {
//		stat_t := fileInfo.Sys().(*syscall.Win32FileAttributeData)
//		return TimestampToTime(stat_t.CreationTime.Nanoseconds()/1e9)
//	}
//
//	// Sys()返回的是interface{}，所以需要类型断言，不同平台需要的类型不一样，linux上为*syscall.Stat_t;
//	// windows 上为 *syscall.Win32FileAttributeData
//	stat_t := fileInfo.Sys().(*syscall.Stat_t)
//	return TimestampToTime(stat_t.Ctim.Sec)
//}

// 获取文件更新时间戳
//func (this *File) GetModifyTimestamp() int64 {
//	fileInfo, _ := os.Stat(this.Name)
//
//	sysType := runtime.GOOS
//	if sysType == "windows" {
//		stat_t := fileInfo.Sys().(*syscall.Win32FileAttributeData)
//		return stat_t.LastWriteTime.Nanoseconds()/1e9
//	}
//
//	stat_t := fileInfo.Sys().(*syscall.Stat_t)
//	return stat_t.Mtim.Sec
//}

// 获取文件更新时间
//func (this *File) GetModifyTime() string {
//	fileInfo, _ := os.Stat(this.Name)
//
//	sysType := runtime.GOOS
//	if sysType == "windows" {
//		stat_t := fileInfo.Sys().(*syscall.Win32FileAttributeData)
//		return TimestampToTime(stat_t.LastWriteTime.Nanoseconds()/1e9)
//	}
//
//	stat_t := fileInfo.Sys().(*syscall.Stat_t)
//	return TimestampToTime(stat_t.Mtim.Sec)
//}

// 获取文件访问时间戳
//func (this *File) GetAccessTimestamp() int64 {
//	fileInfo, _ := os.Stat(this.Name)
//
//	sysType := runtime.GOOS
//	if sysType == "windows" {
//		stat_t := fileInfo.Sys().(*syscall.Win32FileAttributeData)
//		return stat_t.LastAccessTime.Nanoseconds()/1e9
//	}
//
//	stat_t := fileInfo.Sys().(*syscall.Stat_t)
//	return stat_t.Atim.Sec
//}

// 获取文件访问时间
//func (this *File) GetAccessTime() string {
//	fileInfo, _ := os.Stat(this.Name)
//
//	sysType := runtime.GOOS
//	if sysType == "windows" {
//		stat_t := fileInfo.Sys().(*syscall.Win32FileAttributeData)
//		return TimestampToTime(stat_t.LastAccessTime.Nanoseconds()/1e9)
//	}
//
//	stat_t := fileInfo.Sys().(*syscall.Stat_t)
//	return TimestampToTime(stat_t.Atim.Sec)
//}

// 文件是否存在
func (this *File) IsExist() bool {
	if _, err := os.Stat(this.Name); os.IsNotExist(err) {
		return false
	}
	return true
}
