package Config

import (
	"github.com/dengpju/higo-gin/test/app/Services"
)

type Bean struct {

}

func NewBean() *Bean {
	return &Bean{}
}

func (this *Bean)Provider()  {

}

func (this *Bean) DemoService() *Services.DemoService {
	return Services.NewDemoService()
}