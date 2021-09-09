package CoinModel

import (
	"github.com/dengpju/higo-gin/higo"
)


func WithId(v int) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).Id = v
	}
}

func WithUname(v string) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).Uname = v
	}
}

func WithCoin(v int) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).Coin = v
	}
}

