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

var valid *validator.Validate
var ValidContainer map[string]Valid

func init() {
	ValidContainer = make(map[string]Valid)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		valid = v
	} else {
		log.Fatal("error init validator")
	}
}

type IValidate interface {
	RegisterValidator()
}

type Valid map[string]*ValidRules

func (this Valid) Tag(tag string, rules ...*ValidRule) Valid {
	if tag != "" && len(rules) > 0 {
		this[tag] = NewValidRules(rules...)
	}
	return this
}

func RegisterValid(class IClass) Valid {
	v := reflect.ValueOf(class)
	valid := make(Valid, 0)
	ValidContainer[v.Type().String()] = valid
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

func (this *ValidRules) setRule() *ValidRules {
	for _, vrs := range this.Rules {
		this.rule += vrs.Rule + ","
		key := strings.Split(vrs.Rule, "=")
		if len(key) > 1 {
			this.message[key[0]] = vrs.Code
		} else {
			this.message[vrs.Rule] = vrs.Code
		}
	}
	this.rule = strings.Trim(this.rule, ",")
	return this
}

func (this *ValidRules) Rule() string {
	return this.rule
}

func (this *ValidRules) ToFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		v, ok := fl.Field().Interface().(string)
		if ok {
			this.Throw(v)
			return true
		}
		return false
	}
}

func (this *ValidRules) Throw(v interface{}) {
	if err := valid.Var(v, this.rule); err != nil {
		estring := strings.Split(err.Error(), "failed on the '")
		rule := strings.Split(estring[1], "' tag")
		if msg, ok := this.message[rule[0]]; ok {
			panic(NewValidateError(msg))
		}
		panic("validator error")
	}
}

func RegisterValidation(tag string, fn validator.Func) {
	err := valid.RegisterValidation(tag, fn)
	if err != nil {
		panic(fmt.Sprintf("register validator %s error, msg: %s", tag, err.Error()))
	}
}
