package ParamXxxxList

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
)

type XxxxList struct {
}

func New(ctx *gin.Context) *XxxxList {
	param := &XxxxList{}
	higo.Validate(param).Receiver(ctx.ShouldBindJSON(param)).Unwrap()
	return param
}

func (this *XxxxList) RegisterValidator() *higo.Verify {
	return higo.Verifier()
}