package utils

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

// 需要放入token中的信息结构体
type Token struct {
	Payload jwt.MapClaims  // 加密数据
	Secret string  // 秘钥
}

// 构造函数
func NewToken(payload jwt.MapClaims) *Token  {
	secret := os.Getenv("TOKEN_SECRET")
	payload["exp"] = os.Getenv("TOKEN_EXP")
	return &Token{payload,secret}
}

// 创建token
func (this *Token) Create() (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, this.Payload)
	token, err := at.SignedString([]byte(this.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// 解析token
func (this *Token) Parse (token string) (map[string]interface{}, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(this.Secret), nil
	})
	if err != nil {
		return make(map[string]interface{}), err
	}
	return claim.Claims.(jwt.MapClaims), nil
}