package ParamXxxxEdit

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
)

type XxxxEdit struct {
}

func New(ctx *gin.Context) *XxxxEdit {
	param := &XxxxEdit{}
	higo.Validate(param).Receiver(ctx.ShouldBindJSON(param)).Unwrap()
	return param
}

func (this *XxxxEdit) RegisterValidator() *higo.Verify {
	return higo.Verifier()
}