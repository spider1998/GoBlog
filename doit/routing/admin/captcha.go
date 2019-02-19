package admin

import (
	"github.com/go-ozzo/ozzo-routing"
	"Project/doit/service"
	"fmt"
	"Project/doit/form"
	"Project/doit/app"
	"github.com/caeret/ozzo-routing/access"
)

type CaptchaHandler struct{}

func (CaptchaHandler) Generate(c *routing.Context) error {
	token, data, err := service.Captcha.Generate()
	if err != nil {
		return err
	}
	service.SLog.SendLog(form.CreateLogRequest{
		Token:    c.Get(sessionTokenHeaderKey).(string),
		UserType: form.LogUserTypeOperator,
		System:   app.System,
		Action:   "captcha.update",
		IP:       access.GetClientIP(c.Request),
		Remark:   fmt.Sprintf("生成验证码及token： %s。", token),
		Ext:      map[string]interface{}{"token": token},
	})
	return c.Write(map[string]interface{}{
		"token": token,
		"image": data,
	})
}