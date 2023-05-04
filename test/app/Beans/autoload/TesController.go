package autoload

import (
	"fmt"
	"github.com/dunpju/higo-gin/higo"
)

type TestController struct {
}

func NewTestController() *TestController {
	return &TestController{}
}

func (this *TestController) New() higo.IClass {
	return NewTestController()
}

func (this *TestController) Route(hg *higo.Higo) {
	fmt.Println("www")
}
