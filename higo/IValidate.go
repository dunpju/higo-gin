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
	Rules   []*VerifyRule
}

//设置规则
func (this *VerifyRules) setRule() *VerifyRules {
	rules := make([]string, 0)
	for _, vrs := range this.Rules {
		rules = append(rules, vrs.Rule)
		key := strings.Split(vrs.Rule, "=")
		fmt.Println("IValidate:150", vrs, vrs.Rule, key)
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
	fmt.Println(v)
	if msg, ok := v.(code.ICode); ok {
		panic(NewValidateError(msg))
	}
	if err := Validator.Var(v, this.rule); err != nil {
		fmt.Println(err.Error())
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
