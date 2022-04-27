package higo

import (
	"fmt"
	"gitee.com/dengpju/higo-code/code"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
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

//校验者
func Verifier() *Verify {
	return NewVerify()
}

func NewVerify() *Verify {
	return &Verify{Validator: Validator, VerifyRules: make(map[string]*VerifyRules)}
}

type Verify struct {
	Validator      *validator.Validate
	VerifierStruct interface{}
	Verifier       IValidate               //校验者
	VerifyRules    map[string]*VerifyRules //规则map
}

func (this *Verify) Use(validate IValidate) *Verify {
	v := validate.RegisterValidator()
	for tag, r := range v.VerifyRules {
		this.VerifyRules[tag] = r
	}
	return this
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
	return Receiver(values...)
}

//注册规则
func RegisterValidator(validate IValidate) *Verify {
	v := reflect.ValueOf(validate)
	verify := NewVerify()
	verify.Verifier = validate
	VerifyContainer[v.Type().String()] = verify
	return verify
}

//注册自定义校验规则
func RegisterValidation(tag string, fn validator.Func) {
	err := Validator.RegisterValidation(tag, fn)
	if err != nil {
		panic(fmt.Sprintf("register validator %s error, msg: %s", tag, err.Error()))
	}
}

//手动调用注册校验
func Validate(validate IValidate) *Verify {
	verify := validate.RegisterValidator()
	verify.Verifier = validate
	for tag, va := range verify.VerifyRules {
		RegisterValidation(tag, va.ToFunc())
	}
	verify.VerifierStruct = validate
	return verify
}

func Rule(rule string, cod interface{}) *VerifyRule {
	_, ok1 := cod.(code.ICode)
	_, ok2 := cod.(ValidatorToFunc)
	if !ok1 && !ok2 {
		panic(fmt.Errorf("does not support verify rule"))
	}
	return &VerifyRule{Rule: rule, Code: cod}
}

type VerifyRule struct {
	Rule string
	Code interface{}
}

func NewVerifyRules(rules ...*VerifyRule) *VerifyRules {
	vr := &VerifyRules{message: make(map[string]interface{}), Rules: rules}
	return vr.setRule()
}

type VerifyRules struct {
	rule    string
	message map[string]interface{}
	fl      validator.FieldLevel
	Rules   []*VerifyRule
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

type ValidatorToFunc func(fl validator.FieldLevel) (bool, code.ICode)

func (this *VerifyRules) ToFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		this.fl = fl
		if v, ok := fl.Field().Interface().(string); ok {
			this.throw(fl, v)
			return true
		} else if v, ok := fl.Field().Interface().([]string); ok {
			this.throw(fl, v)
			return true
		} else if v, ok := fl.Field().Interface().(int64); ok {
			this.throw(fl, v)
			return true
		} else if v, ok := fl.Field().Interface().([]int64); ok {
			this.throw(fl, v)
			return true
		} else if v, ok := fl.Field().Interface().(uint64); ok {
			this.throw(fl, v)
			return true
		} else if v, ok := fl.Field().Interface().([]uint64); ok {
			this.throw(fl, v)
			return true
		} else if v, ok := fl.Field().Interface().(int); ok {
			this.throw(fl, v)
			return true
		} else if v, ok := fl.Field().Interface().([]int); ok {
			this.throw(fl, v)
			return true
		} else if v, ok := fl.Field().Interface().(float32); ok {
			this.throw(fl, v)
			return true
		} else if v, ok := fl.Field().Interface().([]float32); ok {
			this.throw(fl, v)
			return true
		} else if v, ok := fl.Field().Interface().(float64); ok {
			this.throw(fl, v)
			return true
		} else if v, ok := fl.Field().Interface().([]float64); ok {
			this.throw(fl, v)
			return true
		}
		//未匹配到类型
		for _, msg := range this.message {
			if fn, ok := msg.(ValidatorToFunc); ok {
				boo, v := fn(fl) //自定义函数校验
				if !boo {
					this.throw(fl, v)
					return boo
				}
			}
		}
		return false
	}
}

//抛异常
func (this *VerifyRules) throw(fl validator.FieldLevel, v interface{}) {
	if msg, ok := v.(code.ICode); ok {
		panic(NewValidateError(msg))
	}
	if err := Validator.Var(v, this.rule); err != nil {
		estring := strings.Split(err.Error(), "failed on the '")
		rule := strings.Split(estring[1], "' tag")
		if msg, ok := this.message[rule[0]]; ok {
			if co, ok := msg.(code.ICode); ok {
				panic(NewValidateError(co))
			} else if fn, ok := msg.(ValidatorToFunc); ok {
				boo, v := fn(fl) //自定义函数校验
				if !boo {
					this.throw(fl, v)
				}
				return
			}
		}
		panic("validator error")
	}
}

func (this *VerifyRules) Throw(rule string) {
	keys := strings.Split(rule, "=")
	if len(keys) > 1 {
		rule = keys[0]
	}
	if msg, ok := this.message[rule]; ok {
		if m, ok := msg.(code.ICode); ok {
			panic(NewValidateError(m))
		} else if fn, ok := msg.(ValidatorToFunc); ok {
			boo, v := fn(this.fl) //自定义函数校验
			if !boo {
				this.throw(this.fl, v)
			}
		}
	}
}
