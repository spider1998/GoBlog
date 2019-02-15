package service

import (
	"Project/Doit/app"
	"Project/Doit/code"
	"Project/Doit/entity"
	"Project/Doit/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/utils"
	"github.com/go-ozzo/ozzo-dbx"
	v "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"io/ioutil"
	"net/http"
	"unsafe"
	"mime/multipart"
	"time"
	"crypto/sha1"
	"github.com/gobuffalo/packr/v2/file/resolver/encoding/hex"
)

var User = UserService{
	cachKey:    "cach:us",
	sessionKey: "sess:us:",
	mobileKey:  "mobile:us",
	sessionExp: 3600,
	emailExp:   120,
	mobileExp:  60,
}

type UserService struct {
	sessionKey string
	sessionExp int
	cachKey    string
	emailExp   int
	mobileKey  string
	mobileExp  int
}

//用户注册
func (u *UserService) RegisterUser(request entity.RegisterUserRequest, account string) (user entity.User, err error) {
	err = v.ValidateStruct(&request,
		v.Field(&request.Name, v.Required, v.RuneLength(5, 15)),
		v.Field(&request.Password, v.Required, v.RuneLength(6, 16)),
	)

	vcode, err := app.Redis.Cmd("GET", u.getCachKey(account)).Str()
	if err != nil {
		if err == redis.ErrRespNil {
			err = code.New(http.StatusForbidden, code.CodeUserAccessSessionInvalid).Err("email session not found.")
			return
		}
		err = errors.Wrap(err, "fail to get email code from redis")
		return
	}
	if request.Cach != vcode {
		err = code.New(http.StatusBadRequest, code.CodeVerifyError)
		return
	}

	if err != nil {
		return
	}
	user.ID = uuid.New().String()
	user.Name = request.Name
	user.PasswordHash, err = util.GeneratePasswordHash([]byte(request.Password))
	user.Email = request.Email
	if err != nil {
		err = errors.Wrap(err, "fail to generate password hash")
		return
	}
	user.DatetimeAware = entity.DatetimeAwareNow()
	user.State = entity.UserStateOK

	err = app.DB.Transactional(func(tx *dbx.Tx) error {
		err = tx.Model(&user).Insert()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if util.IsDBDuplicatedErr(err) {
			err = code.New(http.StatusConflict, code.CodeUserExist)
			return
		}
		err = errors.Wrap(err, "fail to register user")
		return
	}
	return
}

//用户登录
func (u *UserService) LoginUser(request entity.LoginUserRequest) (user entity.User, sessionID string, err error) {
	err = v.ValidateStruct(&request,
		v.Field(&request.Name, v.Required),
		v.Field(&request.Password, v.Required),
	)
	if err != nil {
		return
	}
	err = app.DB.Select().Where(dbx.HashExp{"name": request.Name}).One(&user)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusNotFound, code.CodeUserNotExist)
			return
		}
		err = errors.Wrap(err, "fail to find user")
		return
	}

	err = util.ValidatePassword([]byte(request.Password), user.PasswordHash)
	if err != nil {
		err = code.New(http.StatusForbidden, code.CodeUserInvalidPassword)
		return
	}

	sessionID, _ = util.RandString()
	err = app.Redis.Cmd("SET", u.getSessionKey(sessionID), user.ID, "EX", u.sessionExp).Err
	if err != nil {
		err = errors.Wrap(err, "fail to set redis")
		return
	}
	return

}

//检查用户登录状态
func (u *UserService) CheckSession(sessionID string) (user entity.User, err error) {
	userID, err := app.Redis.Cmd("GET", u.getSessionKey(sessionID)).Str()
	if err != nil {
		if err == redis.ErrRespNil {
			err = code.New(http.StatusForbidden, code.CodeUserAccessSessionInvalid).Err("session not found.")
			return
		}
		err = errors.Wrap(err, "fail to get session id from redis")
		return
	}
	err = app.DB.Select().Where(dbx.HashExp{"id": userID}).One(&user)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusNotFound, code.CodeUserAccessSessionInvalid).Err("account not found.")
		}
		err = errors.Wrap(err, "fail to find user")
		return
	}
	err = app.Redis.Cmd("EXPIRE", u.getSessionKey(sessionID), u.sessionExp).Err
	if err != nil {
		err = errors.Wrap(err, "fail to update session expire time")
		return
	}
	return
}

//修改用户资料
func (u *UserService) UpdateInfo(request entity.InfoUpdateRequest) (user entity.User, err error) {
	err = v.ValidateStruct(&request,
		v.Field(&request.Gender, v.Required),
		v.Field(&request.RealName, v.Required, v.RuneLength(5, 15)),
		v.Field(&request.Birthday, v.Required),
		v.Field(&request.Area, v.Required),
		v.Field(&request.HeadImg, v.Required),
		v.Field(&request.Motto, v.Required),
	)
	if err != nil {
		return
	}
	err = app.DB.Select().Where(dbx.HashExp{"id": request.ID}).One(&user)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	user.PersonInfo = request.PersonInfo
	user.AccountInfo = request.AccountInfo
	user.UpdateTime = util.DateTimeStd()
	err = app.DB.Model(&user).Update("Gender", "RealName", "Birthday", "Area", "HeadImg", "Motto", "UpdateTime")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//绑定手机号
func (u *UserService) BindMobile(req entity.BindMobileRequest, uid string) (user entity.User, err error) {
	err = v.ValidateStruct(&req,
		v.Field(&req.Mobile, v.Required),
		v.Field(&req.Mcode, v.Required),
	)
	if err != nil {
		return
	}
	vcode, err := app.Redis.Cmd("GET", u.getMobileKey(uid)).Str()
	if err != nil {
		if err == redis.ErrRespNil {
			err = code.New(http.StatusForbidden, code.CodeUserAccessSessionInvalid).Err("mobile session not found.")
			return
		}
		err = errors.Wrap(err, "fail to get mobile code from redis")
		return
	}
	if req.Mcode != vcode {
		err = code.New(http.StatusBadRequest, code.CodeVerifyError)
		return
	}

	err = app.DB.Select().Where(dbx.HashExp{"id": req.ID}).One(&user)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	user.Mobile = req.Mobile
	user.UpdateTime = util.DateTimeStd()
	err = app.DB.Model(&user).Update("Mobile", "UpdateTime")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return

}

func (u *UserService) SaveAttachment(testH *multipart.FileHeader) (path string) {
	fileDir := testH.Filename + time.Now().String()
	hs := sha1.Sum([]byte(fileDir))
	node := hex.EncodeToString(hs[:])
	path = app.Conf.AttachmentPath + node[:3] + "/"
	return
}

//修改密码
func (u *UserService) SetUserPass(req entity.SetUserPassRequest) (err error) {
	err = v.ValidateStruct(&req,
		v.Field(&req.Password, v.Required, v.RuneLength(6, 16)),
		v.Field(&req.NewPassword, v.Required, v.RuneLength(6, 16)),
	)
	if err != nil {
		return
	}
	var user entity.User
	err = app.DB.Select().Where(dbx.HashExp{"id": req.ID}).One(&user)
	if err != nil {
		if util.IsDBNotFound(err) { //无记录
			err = code.New(http.StatusNotFound, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}

	err = util.ValidatePassword([]byte(req.Password), user.PasswordHash)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = app.DB.Model(&user).Update("PasswordHash")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return

}

//获取邮箱验证码
func (u *UserService) GetVerifyCode(account, req string) (str string, err error) {
	str = xid.New().String()[14:]

	// 创建一个字符串变量，存放相应配置信息
	config :=
		`{"username":"2387805574@qq.com","password":"henuqnarpnucdjci","host":"smtp.qq.com","port":587}`
	// 通过存放配置信息的字符串，创建Email对象
	temail := utils.NewEMail(config)

	temail.To = []string{req}
	temail.From = app.Conf.Email
	temail.Subject = "GoBlog-用户验证"
	temail.HTML = `<html>
		<head>
		</head>
	    	 <body>
			   <div>您的验证码为：` + str + ` </a></div>
	     	</body>
	 	</html>`

	err = temail.Send()
	if err != nil {
		app.Logger.Error().Err(err)
		return
	}
	err = app.Redis.Cmd("SET", u.getCachKey(account), str, "EX", u.emailExp).Err
	if err != nil {
		err = errors.Wrap(err, "fail to set redis")
		return
	}

	return

}

//获取手机验证码
func (u *UserService) GetMobileVerify(uid, mobile string) (code string, err error) {
	var yun entity.YunXun
	vcode := xid.New().String()[14:]
	yun.Sid = app.Conf.Msid
	yun.Token = app.Conf.Mtoken
	yun.Appid = app.Conf.Mappid
	yun.Templateid = app.Conf.Mcach
	yun.Param = vcode + "," + app.Conf.Mexpire
	yun.Mobile = mobile
	yun.Uid = uid
	msg, err := json.Marshal(yun)
	if err != nil {
		fmt.Println(err)
	}
	reader := bytes.NewReader(msg)
	request, err := http.NewRequest("POST", app.Conf.Maddr, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	str := (*string)(unsafe.Pointer(&respBytes))

	var dat map[string]string
	if err = json.Unmarshal([]byte(*str), &dat); err != nil {
		return
	}

	err = app.Redis.Cmd("SET", u.getMobileKey(uid), vcode, "EX", u.mobileExp).Err
	if err != nil {
		err = errors.Wrap(err, "fail to set redis")
		return
	}
	code = dat["msg"]
	return
}


//联系管理员
func (u *UserService) ContactManager(req entity.Contact) (err error) {
	err = v.ValidateStruct(&req,
		v.Field(&req.Name, v.Required),
		v.Field(&req.Mobile, v.Required),
		v.Field(&req.Email, v.Required),
		v.Field(&req.Message, v.Required),
	)
	if err != nil {
		return
	}
	// 创建一个字符串变量，存放相应配置信息
	config :=
		`{"username":"2387805574@qq.com","password":"henuqnarpnucdjci","host":"smtp.qq.com","port":587}`
	// 通过存放配置信息的字符串，创建Email对象
	temail := utils.NewEMail(config)

	temail.To = []string{app.Conf.Email}
	temail.From = app.Conf.Email
	temail.Subject = "GoBlog-用户问题"
	temail.HTML = `<html>
		<head>
		</head>
	    	 <body>
			   <div>用户ID：` + req.UserID + `</div>
			   <div>用户姓名：` + req.Name + `</div>
			   <div>用户姓名：` + req.Name + `</div>
			   <div>用户邮箱：` + req.Email + `</div>
			   <div>反馈问题：` + req.Message + `</div>
	     	</body>
	 	</html>`

	err = temail.Send()
	if err != nil {
		app.Logger.Error().Err(err)
		return
	}
	return
}














func (s *UserService) getSessionKey(sessionID string) string {
	return s.sessionKey + sessionID
}

func (s *UserService) getCachKey(account string) string {
	return s.cachKey + account
}

func (s *UserService) getMobileKey(uid string) string {
	return s.mobileKey + uid
}
