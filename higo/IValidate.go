package higo

import (
	"fmt"
	"gitee.com/dengpju/higo-code/code"
	"github.com/dengpju/higo-throw/exception"
	"github.com/dengpju/higo-utils/utils"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"regexp"
	"strings"
)

var (
	valid          *validator.Validate
	ValidContainer map[string]*Valid
)

func init() {
	ValidContainer = make(map[string]*Valid)
	//初始化校验引擎
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		valid = v
	} else {
		log.Fatal("error init validator")
	}
}

//校验接口
type IValidate interface {
	RegisterValidator() *Valid
}

type Valid struct {
	Verifier   IValidate              //校验者
	ValidRules map[string]*ValidRules //规则map
}

func NewValid() *Valid {
	return &Valid{ValidRules: make(map[string]*ValidRules)}
}

//自定义tag
func (this *Valid) Tag(tag string, rules ...*ValidRule) *Valid {
	if tag != "" && len(rules) > 0 {
		this.ValidRules[tag] = NewValidRules(rules...)
	}
	return this
}

//接收数据
func (this *Valid) Receiver(values ...interface{}) *ErrorResult {
	if len(values) > 0 {
		if err, ok := values[0].(error); ok {
			errStr := exception.ErrorToString(err)
			refVerifierType := reflect.TypeOf(this.Verifier)
			if reflect.Ptr == refVerifierType.Kind() {
				refVerifierType = refVerifierType.Elem()
			}
			binding := ""
			for i := 0; i < refVerifierType.NumField(); i++ {
				reg := regexp.MustCompile("Go struct field " + refVerifierType.Name() + "." + utils.CamelToCase(refVerifierType.Field(i).Name) + " of type") //类型错误
				if reg.MatchString(errStr) {
					binding = utils.CamelToCase(refVerifierType.Field(i).Tag.Get("binding"))
					break
				}
			}
			if "" != binding {
				bindings := strings.Split(binding, ",")
				rules := strings.Split(this.ValidRules[bindings[0]].rule, ",")
				this.ValidRules[bindings[0]].throw(rules[0])//抛出第一规则
			}
		}
	}
	return Receiver(values...)
}

//注册校验规则
func RegisterValid(verifier IValidate) *Valid {
	v := reflect.ValueOf(verifier)
	valid := NewValid()
	valid.Verifier = verifier
	ValidContainer[v.Type().String()] = valid
	return valid
}

//校验者
func Verifier() *Valid {
	return NewValid()
}

//手动调用注册校验
func Validate(verifier IValidate) *Valid {
	valid := verifier.RegisterValidator()
	valid.Verifier = verifier
	for tag, va := range valid.ValidRules {
		RegisterValidation(tag, va.ToFunc())
	}
	return valid
}

type ValidRule struct {
	Rule string
	Code code.ICode
}

func Rule(rule string, code code.ICode) *ValidRule {
	return &ValidRule{Rule: rule, Code: code}
}

type ValidRules struct {
	rule    string
	message map[string]code.ICode
	Rules   []*ValidRule
}

func NewValidRules(rules ...*ValidRule) *ValidRules {
	vr := &ValidRules{message: make(map[string]code.ICode), Rules: rules}
	vr.setRule()
	return vr
}

//设置规则
func (this *ValidRules) setRule() *ValidRules {
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

func (this *ValidRules) Rule() string {
	return this.rule
}

func (this *ValidRules) ToFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		if v, ok := fl.Field().Interface().(string); ok {
			this.Throw(v)
			return true
		} else if v, ok := fl.Field().Interface().(int64); ok {
			this.Throw(v)
			return true
		} else if v, ok := fl.Field().Interface().(int); ok {
			this.Throw(v)
			return true
		}
		return false
	}
}

//抛异常
func (this *ValidRules) Throw(v interface{}) {
	if err := valid.Var(v, this.rule); err != nil {
		estring := strings.Split(err.Error(), "failed on the '")
		rule := strings.Split(estring[1], "' tag")
		this.throw(rule[0])
		panic("validator error")
	}
}

func (this *ValidRules) throw(rule string) {
	if msg, ok := this.message[rule]; ok {
		panic(NewValidateError(msg))
	}
}

//注册自定义校验规则
func RegisterValidation(tag string, fn validator.Func) {
	err := valid.RegisterValidation(tag, fn)
	if err != nil {
		panic(fmt.Sprintf("register validator %s error, msg: %s", tag, err.Error()))
	}
}
