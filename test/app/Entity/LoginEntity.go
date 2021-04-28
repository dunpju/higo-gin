package Entity

type LoginEntity struct {
	UserName    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
	Time        string `json:"time" binding:"required"`
}

func NewLoginEntity() *LoginEntity {
	return &LoginEntity{}
}
