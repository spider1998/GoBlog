package user

import (
	"Project/Doit/handler/session"
	"Project/Doit/handler/user"
	"github.com/go-ozzo/ozzo-routing"
)

func RegisterRoutes(router *routing.RouteGroup) {
	router.Get("/verify/email/<email_account>", user.RegisterVerify) // 拉取邮箱验证码
	router.Get("/verify/mobile/<mobile_account>", user.MobileVerify) // 拉取短信验证码
	router.Post("/register", user.RegisterUser)                      // 注册
	router.Post("/login", user.LoginUser)                            // 登录
	router.Use(user.CheckSession)                                    // 检查登录状态
	router.Get("/sessions/<session_id>", session.GetSession)         // 获取当前登录账户信息
	router.Patch("/information/<user_id>", user.UpdateInfo)          // 修改用户资料
	router.Patch("/info/mobile", user.BindMobile)                    // 绑定手机
	router.Patch("/password/<user_id>", user.SetUserPass)            // 修改密码

}
