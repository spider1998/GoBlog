package admin

import (
	"github.com/go-ozzo/ozzo-routing"
	"strings"
	"Project/doit/code"
	"net/http"
	"Project/doit/service"
	"Project/doit/entity"
)

func AdminRegisterRoutes(router *routing.RouteGroup) {
	var (
		operatorHandler = new(OperatorHandler)
		captchaHandler  = new(CaptchaHandler)
		logHandler      = new(LogHandler)
	)

	{
		router.Get("/captcha", captchaHandler.Generate)							// 拉取图片验证码
		router.Post("/sessions", operatorHandler.SignIn)						// 管理员登录
	}

	router.Use(sessionChecker)

	{
		/*-----------------------------------------Log------------------------------------------------*/
		router.Get("/logs", logHandler.QueryLogs)								// 查询日志
		/*-----------------------------------------User------------------------------------------------*/
		router.Get("/sessions/current", operatorHandler.GetSession)				// 获取管理员信息
		router.Get("/users",operatorHandler.QueryBlogUser)						// 查询用户
		router.Patch("/users/status",operatorHandler.ModifyUserStatus)			// 启用/禁用用户
	}
}


const (
	sessionTokenHeaderKey = "X-Access-Token"
	sessionKey            = "session.operator"
)

func sessionChecker(c *routing.Context) error {
	token := c.Request.Header.Get(sessionTokenHeaderKey)
	if token == "" {
		token = c.Query(strings.ToLower(sessionTokenHeaderKey))
		if token == "" {
			return code.New(http.StatusNotFound, code.CodeOperatorTokenRequired)
		}
	}
	operator, err := service.Operator.CheckToken(token)
	if err != nil {
		return err
	}
	c.Set(sessionKey, operator)
	return c.Next()
}

func getSessionOperator(c *routing.Context) entity.Operator {
	return c.Get(sessionKey).(entity.Operator)
}



