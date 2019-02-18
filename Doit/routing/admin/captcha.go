package admin

import (
	"github.com/go-ozzo/ozzo-routing"
	"Project/Doit/service"
)

type CaptchaHandler struct{}

func (CaptchaHandler) Generate(c *routing.Context) error {
	token, data, err := service.Captcha.Generate()
	if err != nil {
		return err
	}
	return c.Write(map[string]interface{}{
		"token": token,
		"image": data,
	})
}