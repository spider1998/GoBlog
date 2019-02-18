package admin

import (
	"github.com/go-ozzo/ozzo-routing"
)

func AdminRegisterRoutes(router *routing.RouteGroup) {
	var (
		operatorHandler = new(OperatorHandler)
		captchaHandler  = new(CaptchaHandler)
		logHandler      = new(LogHandler)
	)

	{
		router.Get("/captcha", captchaHandler.Generate)
		router.Post("/sessions", operatorHandler.SignIn)
		router.Get("/logs", logHandler.QueryLogs)					//查询日志
	}
}
