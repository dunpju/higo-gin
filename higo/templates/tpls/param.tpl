package {{.Package}}

import (
	"github.com/dunpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
	"sync"
)

var (
	{{.StructName}}Once     sync.Once
	{{.StructName}}Validate *higo.Verify
)

type {{.StructName}} struct {
	//Id    uint64 `form:"id" binding:"id"` // get from the form
	//Id    uint64 `json:"id" binding:"id"` // get from the json
}

func New{{.StructName}}(ctx *gin.Context) *{{.StructName}} {
	param := &{{.StructName}}{}
    {{.StructName}}Once.Do(func() {
    	{{.StructName}}Validate = higo.Validate(param)
    })
	//{{.StructName}}Validate.Receiver(ctx.ShouldBindQuery(param)).Unwrap() // get from the form
	//{{.StructName}}Validate.Receiver(ctx.ShouldBindJSON(param)).Unwrap() // get from the json
	//{{.StructName}}Validate.Receiver(ctx.ShouldBindBodyWith(param, binding.JSON)).Unwrap() // get from the json multiple binding
	return param
}

// RegisterValidator
// https://pkg.go.dev/github.com/go-playground/validator
//
//The custom tag, binding the tag eg: binding:"custom_tag_name"
//require import "gitee.com/dengpju/higo-code/code"
//
//example code:
//func (this *{{.StructName}}) RegisterValidator() *higo.Verify {
//	return higo.RegisterValidator(this).
//		Tag("custom_tag_name1",
//			higo.Rule("required", Codes.Success),
//			higo.Rule("min=5", Codes.Success)).
//		Tag("custom_tag_name2",
//			higo.Rule("required", func() higo.ValidatorToFunc {
//              return func(fl validator.FieldLevel) (bool, code.ICode) {
//                  fmt.Println(fl.Field().Interface())
//                  return true, MinError
//              }
//          }()))
//  Or
//  return higo.Verifier() // Manual call Register Validate: higo.Validate(verifier)
//}
func (this *{{.StructName}}) RegisterValidator() *higo.Verify {
	return higo.Verifier()
	//	.Tag("id",
	//		higo.Rule("required", Codes.Success))
}