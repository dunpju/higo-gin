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
		lowDateTime := t.Elem().FieldByName("CreationTime").Field(0).Uint()
		highDateTime := t.Elem().FieldByName("CreationTime").Field(1).Uint()
		return Nanoseconds(uint32(lowDateTime), uint32(highDateTime)) / 1e9
	}
	return int64(t.Elem().FieldByName("Ctim").FieldByName("Sec").Int())
}

// 获取文件创建时间
func (this *File) GetCreateTime() string {
	fileInfo, _ := os.Stat(this.Name)
	t := reflect.ValueOf(fileInfo.Sys())
	var sec int64
	if "syscall.Win32FileAttributeData" == fmt.Sprintf("%s", t.Elem().Type()) {
		lowDateTime := t.Elem().FieldByName("CreationTime").Field(0).Uint()
		highDateTime := t.Elem().FieldByName("CreationTime").Field(1).Uint()
		sec = Nanoseconds(uint32(lowDateTime), uint32(highDateTime)) / 1e9
	}else {
		sec = int64(t.Elem().FieldByName("Ctim").FieldByName("Sec").Int())
	}
	return TimestampToTime(sec)
}

// 获取文件更新时间戳
func (this *File) GetModifyTimestamp() int64 {
	fileInfo, _ := os.Stat(this.Name)
	t := reflect.ValueOf(fileInfo.Sys())
	if "syscall.Win32FileAttributeData" == fmt.Sprintf("%s", t.Elem().Type()) {
		lowDateTime := t.Elem().FieldByName("LastWriteTime").Field(0).Uint()
		highDateTime := t.Elem().FieldByName("LastWriteTime").Field(1).Uint()
		return Nanoseconds(uint32(lowDateTime), uint32(highDateTime)) / 1e9
	}
	return int64(t.Elem().FieldByName("Mtim").FieldByName("Sec").Int())
}

// 获取文件更新时间
func (this *File) GetModifyTime() string {
	fileInfo, _ := os.Stat(this.Name)
	t := reflect.ValueOf(fileInfo.Sys())
	var sec int64
	if "syscall.Win32FileAttributeData" == fmt.Sprintf("%s", t.Elem().Type()) {
		lowDateTime := t.Elem().FieldByName("LastWriteTime").Field(0).Uint()
		highDateTime := t.Elem().FieldByName("LastWriteTime").Field(1).Uint()
		sec = Nanoseconds(uint32(lowDateTime), uint32(highDateTime)) / 1e9
	} else {
		sec = int64(t.Elem().FieldByName("Mtim").FieldByName("Sec").Int())
	}
	return TimestampToTime(sec)
}

// 获取文件访问时间戳
func (this *File) GetAccessTimestamp() int64 {
	fileInfo, _ := os.Stat(this.Name)
	t := reflect.ValueOf(fileInfo.Sys())
	if "syscall.Win32FileAttributeData" == fmt.Sprintf("%s", t.Elem().Type()) {
		lowDateTime := t.Elem().FieldByName("LastAccessTime").Field(0).Uint()
		highDateTime := t.Elem().FieldByName("LastAccessTime").Field(1).Uint()
		return Nanoseconds(uint32(lowDateTime), uint32(highDateTime)) / 1e9
	}
	return int64(t.Elem().FieldByName("Atim").FieldByName("Sec").Int())
}

// 获取文件访问时间
func (this *File) GetAccessTime() string {
	fileInfo, _ := os.Stat(this.Name)
	t := reflect.ValueOf(fileInfo.Sys())
	var sec int64
	if "syscall.Win32FileAttributeData" == fmt.Sprintf("%s", t.Elem().Type()) {
		lowDateTime := t.Elem().FieldByName("LastAccessTime").Field(0).Uint()
		highDateTime := t.Elem().FieldByName("LastAccessTime").Field(1).Uint()
		sec = Nanoseconds(uint32(lowDateTime), uint32(highDateTime)) / 1e9
	}else {
		sec = int64(t.Elem().FieldByName("Atim").FieldByName("Sec").Int())
	}
	return TimestampToTime(sec)
}

// 文件是否存在
func (this *File) IsExist() bool {
	if _, err := os.Stat(this.Name); os.IsNotExist(err) {
		return false
	}
	return true
}
