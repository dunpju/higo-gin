package EnumSearchType

import "github.com/dengpju/higo-enum/enum"

var e SearchType

func Inspect(value int) error {
	return e.Inspect(value)
}

//搜索类型
type SearchType int

func (this SearchType) Name() string {
	return "SearchType"
}

func (this SearchType) Inspect(value interface{}) error {
	return enum.Inspect(this, value)
}

func (this SearchType) Message() string {
	return enum.String(this)
}

const (
	Anecdote SearchType = 1 //轶事记录
	DomainNorm SearchType = 2 //领域指标
	ChildName SearchType = 3 //幼儿
	Label SearchType = 4 //标签
)

func (this SearchType) Register() enum.Message {
	return make(enum.Message).
	    Put(Anecdote, "轶事记录").
	    Put(DomainNorm, "领域指标").
	    Put(ChildName, "幼儿").
	    Put(Label, "标签")
}