package utils

import (
	"fmt"
	"os"
	"reflect"
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
	t := reflect.ValueOf(fileInfo.Sys())
	if "syscall.Win32FileAttributeData" == fmt.Sprintf("%s", t.Elem().Type()) {
		lowDateTime := t.Elem().FieldByName("CreationTime").FieldByName("LowDateTime").Uint()
		highDateTime := t.Elem().FieldByName("CreationTime").FieldByName("HighDateTime").Uint()
		return Nanoseconds(uint32(lowDateTime), uint32(highDateTime)) / 1e9
	}
	return int64(t.Elem().FieldByName("Ctim").FieldByName("Sec").Int())
}

// 获取文件更新时间戳
func (this *File) GetModifyTimestamp() int64 {
	fileInfo, _ := os.Stat(this.Name)
	t := reflect.ValueOf(fileInfo.Sys())
	if "syscall.Win32FileAttributeData" == fmt.Sprintf("%s", t.Elem().Type()) {
		lowDateTime := t.Elem().FieldByName("LastWriteTime").FieldByName("LowDateTime").Uint()
		highDateTime := t.Elem().FieldByName("LastWriteTime").FieldByName("HighDateTime").Uint()
		return Nanoseconds(uint32(lowDateTime), uint32(highDateTime)) / 1e9
	}
	return int64(t.Elem().FieldByName("Mtim").FieldByName("Sec").Int())
}

// 获取文件访问时间戳
func (this *File) GetAccessTimestamp() int64 {
	fileInfo, _ := os.Stat(this.Name)
	t := reflect.ValueOf(fileInfo.Sys())
	if "syscall.Win32FileAttributeData" == fmt.Sprintf("%s", t.Elem().Type()) {
		lowDateTime := t.Elem().FieldByName("LastAccessTime").FieldByName("LowDateTime").Uint()
		highDateTime := t.Elem().FieldByName("LastAccessTime").FieldByName("HighDateTime").Uint()
		return Nanoseconds(uint32(lowDateTime), uint32(highDateTime)) / 1e9
	}
	return int64(t.Elem().FieldByName("Atim").FieldByName("Sec").Int())
}

// 获取文件创建时间
func (this *File) GetCreateTime() string {
	timestamp := this.GetCreateTimestamp()
	return TimestampToTime(timestamp)
}

// 获取文件更新时间
func (this *File) GetModifyTime() string {
	timestamp := this.GetModifyTimestamp()
	return TimestampToTime(timestamp)
}

// 获取文件访问时间
func (this *File) GetAccessTime() string {
	timestamp := this.GetAccessTimestamp()
	return TimestampToTime(timestamp)
}

// 文件是否存在
func (this *File) IsExist() bool {
	if _, err := os.Stat(this.Name); os.IsNotExist(err) {
		return false
	}
	return true
}
