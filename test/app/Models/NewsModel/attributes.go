package NewsModel

import (
	"github.com/dunpju/higo-gin/higo"
)

const (
	NewsId     = "news_id"     //主键
	Title      = "title"       //标题
	Clicknum   = "clicknum"    //点击量
	CreateTime = "create_time" //创建时间
)

func WithNewsId(v int) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).NewsId = v
	}
}

func WithTitle(v string) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).Title = v
	}
}

func WithClicknum(v int) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).Clicknum = v
	}
}

func WithCreateTime(v interface{}) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).CreateTime = v
	}
}
