package AdminModel

import (
	"github.com/dunpju/higo-gin/higo"
	"time"
)

const (
	AdminId   = "admin_id"
	AdminName = "admin_name"
	UserId    = "user_id"
)

func WithAdminId(v int64) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).AdminId = v
	}
}

func WithAdminName(v string) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).AdminName = v
	}
}

func WithUserId(v int64) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).UserId = v
	}
}

func WithState(v int) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).State = v
	}
}

func WithIsSuper(v int) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).IsSuper = v
	}
}

func WithPassword(v string) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).Password = v
	}
}

func WithCreateTime(v time.Time) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).CreateTime = v
	}
}

func WithUpdateTime(v time.Time) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).UpdateTime = v
	}
}

func WithDeleteTime(v interface{}) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).DeleteTime = v
	}
}
