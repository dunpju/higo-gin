package higo

import (
	"fmt"
	"gitee.com/dengpju/higo-code/code"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"strings"
	"sync"
)

var (
	Validator       *validator.Validate
	VerifyContainer *sync.Map
)

func init() {
	VerifyContainer = &sync.Map{}
	//初始化校验引擎
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		Validator = v
	} else {
		log.Fatal("error init validator")
	}
}

// IValidate 校验接口
type IValidate interface {
	RegisterValidator() *Verify
}

// Verifier 校验者
func Verifier() *Verify {
	return NewVerify()
}

func NewVerify() *Verify {
	return &Verify{Validator: Validator, VerifyRules: &sync.Map{}}
}

type Verify struct {
	Validator      *validator.Validate
	VerifierStruct interface{}
	Verifier       IValidate //校验者
	VerifyRules    *sync.Map //规则map map[string]*VerifyRules
}

func (this *Verify) Use(validate IValidate, validates ...IValidate) *Verify {
	v := validate.RegisterValidator()
	v.VerifyRules.Range(func(tag, r interface{}) bool {
		this.VerifyRules.Store(tag, r)
		return true
	})
	for _, validate := range validates {
		v := validate.RegisterValidator()
		v.VerifyRules.Range(func(tag, r interface{}) bool {
			this.VerifyRules.Store(tag, r)
			return true
		})
	}
	return this
}

// Tag 自定义tag
func (this *Verify) Tag(tag string, rules ...*VerifyRule) *Verify {
	if tag != "" && len(rules) > 0 {
		this.VerifyRules.Store(tag, NewRuleGroup(tag, rules...))
	}
	return this
}

// Struct 服务层校验 struct https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
func (this *Verify) Struct() *ErrorResult {
	return Receiver(this.Validator.Struct(this.VerifierStruct))
}

func (this *Verify) Unwrap() interface{} {
	return this.Struct().Unwrap()
}

// Receiver 接收数据
func (this *Verify) Receiver(values ...interface{}) *ErrorResult {
	return Receiver(values...)
}

// RegisterValidator 注册规则
func RegisterValidator(validate IValidate) *Verify {
	v := reflect.ValueOf(validate)
	verify := NewVerify()
	verify.Verifier = validate
	VerifyContainer.Store(v.Type().String(), verify)
	return verify
}

// RegisterValidation 注册自定义校验规则
func RegisterValidation(tag string, fn validator.Func) {
	err := Validator.RegisterValidation(tag, fn)
	if err != nil {
		panic(fmt.Sprintf("register validator %s error, msg: %s", tag, err.Error()))
	}
}

// Validate 手动调用注册校验
func Validate(validate IValidate) *Verify {
	verify := validate.RegisterValidator()
	verify.Verifier = validate
	verify.VerifyRules.Range(func(tag, va interface{}) bool {
		RegisterValidation(tag.(string), va.(*RuleGroup).ToFunc())
		return true
	})
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

func RuleFunc(rule string, cod ValidatorToFunc) *VerifyRule {
	return &VerifyRule{Rule: rule, Code: cod}
}

type VerifyRule struct {
	Rule string
	Code interface{}
}

func NewRuleGroup(tag string, rules ...*VerifyRule) *RuleGroup {
	rg := &RuleGroup{tag: tag, message: &sync.Map{}, Rules: rules}
	return rg.setRule()
}

// RuleGroup 规则组
type RuleGroup struct {
	tag     string
	rule    string
	message *sync.Map //map[string]interface{}
	fl      validator.FieldLevel
	Rules   []*VerifyRule
}

// 设置规则
func (this *RuleGroup) setRule() *RuleGroup {
	rules := make([]string, 0)
	for _, vrs := range this.Rules {
		rules = append(rules, vrs.Rule)
		key := strings.Split(vrs.Rule, "=")
		if len(key) > 1 {
			this.message.Store(key[0], vrs.Code)
		} else {
			this.message.Store(vrs.Rule, vrs.Code)
		}
	}
	// 转换 required,min=4
	this.rule = strings.Join(rules, ",")
	return this
}

func (this *RuleGroup) Rule() string {
	return this.rule
}

type ValidatorToFunc func(fieldLevel validator.FieldLevel) (bool, code.ICode)

// ToFunc
// 自定义规则校验顺序优先与gin原始binging tag `json:"xxx" binding:"required"`
// 自定义tag转换Func
func (this *RuleGroup) ToFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		this.fl = fl
		// 遍历规则
		for _, rule := range this.Rules {
			if customValidFunc, ok := rule.Code.(ValidatorToFunc); ok {
				b, c := customValidFunc(fl)
				if !b {
					this.valid(rule.Rule, fl, c)
				}
			} else {
				this.valid(rule.Rule, fl, fl.Field().Interface())
			}
		}
		return true
	}
}

// 校验值
func (this *RuleGroup) valid(rule string, fl validator.FieldLevel, v interface{}) {
	if msg, ok := v.(code.ICode); ok {
		panic(NewValidateError(msg))
	}
	if err := Validator.Var(v, this.rule); err != nil {
		key := strings.Split(rule, "=")
		if len(key) > 1 {
			rule = key[0]
		}
		if msg, ok := this.message.Load(rule); ok {
			if co, ok := msg.(code.ICode); ok {
				panic(NewValidateError(co))
			} else if fn, ok := msg.(ValidatorToFunc); ok {
				boo, v := fn(fl) //自定义函数校验
				if !boo {
					this.valid(rule, fl, v)
				}
				return
			}
		}
		panic("validator error")
	}
}

// 抛异常
func (this *RuleGroup) Throw(rule string) {
	keys := strings.Split(rule, "=")
	if len(keys) > 1 {
		rule = keys[0]
	}
	if msg, ok := this.message.Load(rule); ok {
		if m, ok := msg.(code.ICode); ok {
			panic(NewValidateError(m))
		} else if fn, ok := msg.(ValidatorToFunc); ok {
			boo, v := fn(this.fl) //自定义函数校验
			if !boo {
				this.valid(rule, this.fl, v)
			}
		}
	}
}
