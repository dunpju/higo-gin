package higo

import (
	"fmt"
	"gitee.com/dengpju/higo-code/code"
	"github.com/dengpju/higo-throw/exception"
	"github.com/dengpju/higo-utils/utils/stringutil"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"regexp"
	"strings"
)

var (
	Validator       *validator.Validate
	VerifyContainer map[string]*Verify
)

func init() {
	VerifyContainer = make(map[string]*Verify)
	//初始化校验引擎
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		Validator = v
	} else {
		log.Fatal("error init validator")
	}
}

//校验接口
type IValidate interface {
	RegisterValidator() *Verify
}

type Verify struct {
	Validator      *validator.Validate
	VerifierStruct interface{}
	Verifier       IValidate               //校验者
	VerifyRules    map[string]*VerifyRules //规则map
}

func NewVerify() *Verify {
	return &Verify{Validator: Validator, VerifyRules: make(map[string]*VerifyRules)}
}

//自定义tag
func (this *Verify) Tag(tag string, rules ...*VerifyRule) *Verify {
	if tag != "" && len(rules) > 0 {
		this.VerifyRules[tag] = NewVerifyRules(rules...)
	}
	return this
}

//服务层校验 struct https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
func (this *Verify) Struct() *ErrorResult {
	return Receiver(this.Validator.Struct(this.VerifierStruct))
}

func (this *Verify) Unwrap() interface{} {
	return this.Struct().Unwrap()
}

//接收数据
func (this *Verify) Receiver(values ...interface{}) *ErrorResult {
	if len(values) > 0 {
		if err, ok := values[0].(error); ok {
			errStr := exception.ErrorToString(err)
			refVerifierType := reflect.TypeOf(this.Verifier)
			if reflect.Ptr == refVerifierType.Kind() {
				refVerifierType = refVerifierType.Elem()
			}
			binding := ""
			for i := 0; i < refVerifierType.NumField(); i++ {
				reg := regexp.MustCompile("Go struct field " + refVerifierType.Name() + "." + stringutil.CamelToCase(refVerifierType.Field(i).Name) + " of type") //类型错误
				if reg.MatchString(errStr) {
					binding = stringutil.CamelToCase(refVerifierType.Field(i).Tag.Get("binding"))
					break
				}
			}
			if "" != binding {
				bindings := strings.Split(binding, ",")
				rules := strings.Split(this.VerifyRules[bindings[0]].rule, ",")
				this.VerifyRules[bindings[0]].throw(rules[0]) //抛出第一规则
			}
		}
	}
	return Receiver(values...)
}

//注册校验规则
func RegisterValidator(verifier IValidate) *Verify {
	v := reflect.ValueOf(verifier)
	verify := NewVerify()
	verify.Verifier = verifier
	VerifyContainer[v.Type().String()] = verify
	return verify
}

//校验者
func Verifier() *Verify {
	return NewVerify()
}

//手动调用注册校验
func Validate(verifier IValidate) *Verify {
	verify := verifier.RegisterValidator()
	verify.Verifier = verifier
	for tag, va := range verify.VerifyRules {
		RegisterValidation(tag, va.ToFunc())
	}
	verify.VerifierStruct = verifier
	return verify
}

type VerifyRule struct {
	Rule string
	Code code.ICode
}

func Rule(rule string, code code.ICode) *VerifyRule {
	return &VerifyRule{Rule: rule, Code: code}
}

type VerifyRules struct {
	rule    string
	message map[string]code.ICode
	Rules   []*VerifyRule
}

func NewVerifyRules(rules ...*VerifyRule) *VerifyRules {
	vr := &VerifyRules{message: make(map[string]code.ICode), Rules: rules}
	vr.setRule()
	return vr
}

//设置规则
func (this *VerifyRules) setRule() *VerifyRules {
	rules := make([]string, 0)
	for _, vrs := range this.Rules {
		rules = append(rules, vrs.Rule)
		key := strings.Split(vrs.Rule, "=")
		if len(key) > 1 {
			this.message[key[0]] = vrs.Code
		} else {
			this.message[vrs.Rule] = vrs.Code
		}
	}
	this.rule = strings.Join(rules, ",")
	return this
}

func (this *VerifyRules) Rule() string {
	return this.rule
}

func (this *VerifyRules) ToFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		if v, ok := fl.Field().Interface().(string); ok {
			this.Throw(v)
			return true
		} else if v, ok := fl.Field().Interface().([]string); ok {
			this.Throw(v)
			return true
		} else if v, ok := fl.Field().Interface().(int64); ok {
			this.Throw(v)
			return true
		} else if v, ok := fl.Field().Interface().([]int64); ok {
			this.Throw(v)
			return true
		} else if v, ok := fl.Field().Interface().(int); ok {
			this.Throw(v)
			return true
		} else if v, ok := fl.Field().Interface().([]int); ok {
			this.Throw(v)
			return true
		} else if v, ok := fl.Field().Interface().(float32); ok {
			this.Throw(v)
			return true
		} else if v, ok := fl.Field().Interface().([]float32); ok {
			this.Throw(v)
			return true
		} else if v, ok := fl.Field().Interface().(float64); ok {
			this.Throw(v)
			return true
		} else if v, ok := fl.Field().Interface().([]float64); ok {
			this.Throw(v)
			return true
		}
		return false
	}
}

//抛异常
func (this *VerifyRules) Throw(v interface{}) {
	if err := Validator.Var(v, this.rule); err != nil {
		estring := strings.Split(err.Error(), "failed on the '")
		rule := strings.Split(estring[1], "' tag")
		this.throw(rule[0])
		panic("validator error")
	}
}

func (this *VerifyRules) throw(rule string) {
	if msg, ok := this.message[rule]; ok {
		panic(NewValidateError(msg))
	}
}

//注册自定义校验规则
func RegisterValidation(tag string, fn validator.Func) {
	err := Validator.RegisterValidation(tag, fn)
	if err != nil {
		panic(fmt.Sprintf("register validator %s error, msg: %s", tag, err.Error()))
	}
}
