package {{.Package}}

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
)

type {{.StructName}} struct {
}

func New(ctx *gin.Context) *{{.StructName}} {
	param := &{{.StructName}}{}
	higo.Validate(param).Receiver(ctx.ShouldBindJSON(param)).Unwrap()
	return param
}

func (this *{{.StructName}}) RegisterValidator() *higo.Valid {
	return higo.Verifier()
}