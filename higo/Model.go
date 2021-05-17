package higo

import (
	"encoding/json"
)

type Model interface {
	New() IClass
	Mutate(attrs ...Property) Model
}

type Models string

func MakeModels(v interface{}) Models {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return Models(b)
}
