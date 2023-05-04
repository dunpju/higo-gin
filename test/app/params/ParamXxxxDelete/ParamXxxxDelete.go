package ParamXxxxDelete

import (
	"github.com/dunpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
)

type XxxxDelete struct {
}

func New(ctx *gin.Context) *XxxxDelete {
	param := &XxxxDelete{}
	higo.Validate(param).Receiver(ctx.ShouldBindJSON(param)).Unwrap()
	return param
}

func (this *XxxxDelete) RegisterValidator() *higo.Verify {
	return higo.Verifier()
}
