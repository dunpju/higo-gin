package EnumSearchType

import (
	"fmt"
)

var (
	enums map[SearchType]*enum
)

const (
	Anecdote SearchType = 1 //轶事记录
	DomainNorm SearchType = 2 //领域指标
	ChildName SearchType = 3 //幼儿
	Label SearchType = 4 //标签
)

func init() {
	enums = make(map[SearchType]*enum)
	enums[Anecdote] = newEnum(int(Anecdote), "轶事记录")
	enums[DomainNorm] = newEnum(int(DomainNorm), "领域指标")
	enums[ChildName] = newEnum(int(ChildName), "幼儿")
	enums[Label] = newEnum(int(Label), "标签")
}

type enum struct {
	code    int
	message string
}

func newEnum(code int, message string) *enum {
	return &enum{code: code, message: message}
}

func Enums() map[SearchType]*enum {
	return enums
}

func Inspect(value int) error {
	_, err := SearchType(value).inspect()
	if err != nil {
		return err
	}
	return nil
}

// SearchType 搜索类型
type SearchType int

func (this SearchType) inspect() (*enum, error) {
	if e, ok := enums[this]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%d enum undefined", this)
}

func (this SearchType) get() *enum {
	e, err := this.inspect()
	if err != nil {
		panic(err)
	}
	return e
}

func (this SearchType) Code() int {
	return this.get().code
}

func (this SearchType) Message() string {
	return this.get().message
}
