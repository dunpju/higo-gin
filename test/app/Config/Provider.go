package Config

import (
	"github.com/dengpju/higo-gin/test/app/Services"
)

type Provider struct {

}

func NewProvider() *Provider {
	return &Provider{}
}

func (this *Provider) DemoService() *Services.DemoService {
	return Services.NewDemoService()
}