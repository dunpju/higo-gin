package UserModel

import "github.com/dengpju/higo-gin/higo"

func WithId(id int) higo.Property {
	return func(class higo.IClass) {
		class.(*UserModelImpl).Id = id
	}
}

func WithUname(name string) higo.Property {
	return func(class higo.IClass) {
		class.(*UserModelImpl).Uname = name
	}
}

func WithUtel(tel string) higo.Property {
	return func(class higo.IClass) {
		class.(*UserModelImpl).Utel = tel
	}
}
