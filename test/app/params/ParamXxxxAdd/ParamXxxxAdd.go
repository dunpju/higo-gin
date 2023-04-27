package ParamXxxxAdd

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
)

type XxxxAdd struct {
}

func New(ctx *gin.Context) *XxxxAdd {
	param := &XxxxAdd{}
	higo.Validate(param).Receiver(ctx.ShouldBindJSON(param)).Unwrap()
	return param
}

func (this *XxxxAdd) RegisterValidator() *higo.Verify {
	return higo.Verifier()
}