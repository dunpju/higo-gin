package V3

import (
	"github.com/dengpju/higo-gin/higo"
)

func init() {
	higo.Once.Do(func() {
		higo.AddContainer(&DemoController{})
	})
}
