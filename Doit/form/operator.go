package form

type CreateOperatorRequest struct {
	Name        string            `json:"name"`
	Password    string            `json:"password"`
	RealName    string            `json:"real_name"`
	Permissions map[string]string `json:"permissions"`
}

type UpdateOperatorRequest struct {
	ID          string            `json:"-"`
	RealName    string            `json:"real_name"`
	Permissions map[string]string `json:"permissions"`
}

type OperatorSignInRequest struct {
	Name         string `json:"name"`
	Password     string `json:"password"`
	CaptchaToken string `json:"captcha_token"`
	CaptchaCode  string `json:"captcha_code"`
}

type SetOperatorShortcutsRequest struct {
	ID        string        `json:"-"`
	Shortcuts []interface{} `json:"shortcuts"`
}
