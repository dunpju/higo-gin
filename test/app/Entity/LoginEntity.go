package Entity

type LoginEntity struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	CaptchaCode string `json:"captcha_code"`
	Time        string `json:"time"`
}

func NewLoginEntity() *LoginEntity {
	return &LoginEntity{}
}
