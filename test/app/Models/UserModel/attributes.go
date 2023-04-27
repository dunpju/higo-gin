package UserModel

import (
	"github.com/dengpju/higo-gin/higo"
)

const (
    Id = "id"  //
    Uname = "uname"  //
    UTel = "u_tel"  //
    Score = "score"  //
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

func WithUTel(v string) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).UTel = v
	}
}

func WithScore(v int) higo.Property {
	return func(class higo.IClass) {
		class.(*Impl).Score = v
	}
}

