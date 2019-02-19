package routing

import (
	"Project/doit/handler/session"
	"Project/doit/handler/user"
	"github.com/go-ozzo/ozzo-routing"
	"fmt"
)

func UserRegisterRoutes(router *routing.RouteGroup) {
	router.Post("/upload", func(context *routing.Context) error {
		fmt.Println("test upload file")
		return nil
	},user.AttachUpload)
	router.Get("/verify/email/<email_account>", user.RegisterVerify) // 拉取邮箱验证码
	router.Get("/verify/mobile/<mobile_account>", user.MobileVerify) // 拉取短信验证码
	router.Post("/register", user.RegisterUser)                      // 注册
	router.Post("/login", user.LoginUser)                            // 登录
	router.Use(user.CheckSession)                                    	  // 检查登录状态
	router.Get("/sessions/<session_id>", session.GetSession)         // 获取当前登录账户信息
	router.Patch("/information/<user_id>", user.UpdateInfo)          // 修改用户资料
	router.Patch("/info/mobile", user.BindMobile)                    // 绑定手机
	router.Patch("/password/<user_id>", user.SetUserPass)            // 修改密码
	router.Post("/contact",user.Contact)							  // 联系管理员


}
