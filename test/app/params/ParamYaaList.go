package params

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/gin-gonic/gin"
)

type YaaList struct {
	//Id    uint64 `form:"id" binding:"id"` // get from the form
	//Id    uint64 `json:"id" binding:"id"` // get from the json
}

func NewYaaList(ctx *gin.Context) *YaaList {
	param := &YaaList{}
	//higo.Validate(param).Receiver(ctx.ShouldBindQuery(param)).Unwrap() // get from the form
	//higo.Validate(param).Receiver(ctx.ShouldBindJSON(param)).Unwrap() // get from the json
	//higo.Validate(param).Receiver(ctx.ShouldBindBodyWith(param, binding.JSON)).Unwrap() // get from the json multiple binding
	return param
}

//The custom tag, binding the tag eg: binding:"custom_tag_name"
//require import "gitee.com/dengpju/higo-code/code"
//
//example code:
//func (this *YaaList) RegisterValidator() *higo.Verify {
//	return higo.RegisterValidator(this).
//		Tag("custom_tag_name",
//			higo.Rule("required", Codes.Success),
//			higo.Rule("min=5", Codes.Success))
//  Or
//  return higo.Verifier() // Manual call Register Validate: higo.Validate(verifier)
//}
func (this *YaaList) RegisterValidator() *higo.Verify {
	return higo.Verifier()
	//	.Tag("id",
	//		higo.Rule("required", Codes.Success))
}