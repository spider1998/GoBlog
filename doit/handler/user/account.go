package user

import (
	"Project/doit/code"
	"Project/doit/entity"
	"Project/doit/handler/session"
	"Project/doit/service"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

/*---------------------------------获取验证码-----------------------------------------------------------------*/

//获取邮箱验证码
func RegisterVerify(c *routing.Context) error {
	account := access.GetClientIP(c.Request)
	email := c.Param("email_account")
	cach, err := service.User.GetVerifyCode(account, email)
	if err != nil {
		return err
	}
	return c.Write(len(cach))
}

//获取手机验证码
func MobileVerify(c *routing.Context) error {
	client := access.GetClientIP(c.Request)
	uid := strings.Replace(client, ".", "", -1)
	mobile := c.Param("mobile_account")
	statuCode, err := service.User.GetMobileVerify(uid, mobile)
	if err != nil {
		return err
	}
	return c.Write(statuCode)
}

/*---------------------------------用户层操作-----------------------------------------------------------------*/

//用户注册
func RegisterUser(c *routing.Context) error {
	var req entity.RegisterUserRequest
	account := access.GetClientIP(c.Request)
	err := c.Read(&req)
	if err != nil {
		return code.New(400, code.CodeBadRequest).Err(err)
	}

	user, err := service.User.RegisterUser(req, account)
	if err != nil {
		return err
	}
	return c.Write(user)

}

//用户登录
func LoginUser(c *routing.Context) error {
	var req entity.LoginUserRequest
	err := c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	_, sessionID, err := service.User.LoginUser(req)
	if err != nil {
		return err
	}
	c.Response.Header().Set("X-Access-Session", sessionID) //写入响应头
	return c.Write(sessionID)
}

//检查用户状态
func CheckSession(c *routing.Context) error {
	accessSession := c.Request.Header.Get("X-Access-Session")
	if accessSession == "" {
		accessSession = c.Param("session_id")
		if accessSession == "" {
			return code.New(http.StatusForbidden, code.CodeUserAccessSessionInvalid).Err("Session not found.")
		}
	}
	user, err := service.User.CheckSession(accessSession)
	if err != nil {
		return err
	}
	if user.State != entity.UserStateOK {
		return code.New(http.StatusForbidden, code.CodeUserAccessSessionInvalid).Err("user status is invalid.")
	}
	session.SetUseression(c, user)
	return c.Next()
}

//修改用户资料
func UpdateInfo(c *routing.Context) error {
	var req entity.InfoUpdateRequest
	err := c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}

	req.ID = c.Param("user_id")
	newInfo, err := service.User.UpdateInfo(req)
	if err != nil {
		return errors.Wrap(err, "fail to edit user information")
	}
	return c.Write(newInfo)
}

//绑定手机号
func BindMobile(c *routing.Context) (err error) {
	var req entity.BindMobileRequest
	err = c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	client := access.GetClientIP(c.Request)
	uid := strings.Replace(client, ".", "", -1)
	user, err := service.User.BindMobile(req, uid)
	if err != nil {
		return
	}
	return c.Write(user)
}

//修改密码
func SetUserPass(c *routing.Context) (err error) {
	var req entity.SetUserPassRequest
	err = c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}

	req.ID = c.Param("user_id")
	err = service.User.SetUserPass(req)
	if err != nil {
		return
	}
	return
}

func AttachUpload(c *routing.Context) (err error) {
	return
}

//联系管理员
func Contact(c *routing.Context) error {
	var contact entity.Contact
	err := c.Read(&contact)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	userID := session.GetUserSession(c).ID
	contact.UserID = userID
	err = service.User.ContactManager(contact)
	if err != nil {
		return err
	}
	return err
}


