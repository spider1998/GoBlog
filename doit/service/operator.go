package service

import (
	"Project/doit/form"
	"Project/doit/entity"
	"Project/doit/app"
	"Project/doit/code"
	"Project/doit/util"

	v "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/mediocregopher/radix.v2/redis"
)

var Operator = &OperatorService{}

type OperatorService struct{}

//管理员登陆
func (s *OperatorService) SignIn(request form.OperatorSignInRequest) (token string, operator entity.Operator, err error) {
	err = v.ValidateStruct(&request,
		v.Field(&request.Name, v.Required),
		v.Field(&request.Password, v.Required),
		v.Field(&request.CaptchaToken, v.Required),
		v.Field(&request.CaptchaCode, v.Required),
	)
	if err != nil {
		return
	}

	err = Captcha.Validate(request.CaptchaToken, request.CaptchaCode)
	if err != nil {
		return
	}

	err = app.DB.Select().Where(dbx.HashExp{"name":request.Name}).One(&operator)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusNotFound, code.CodeUserNotExist)
			return
		}
		err = errors.Wrap(err, "fail to find user")
		return
	}
	if operator.State != entity.OperatorStateEnabled {
		err = code.New(http.StatusBadRequest,code.CodeUserDisabled)
		return
	}
	err = util.ValidatePassword([]byte(request.Password), operator.PasswordHash)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			err = code.New(http.StatusBadRequest,code.CodeUserInvalidPassword)
			return
		}
		err = errors.WithStack(err)
		return
	}

	token = RandString(32)
	err = app.Redis.Cmd("SET", "go-blog:op:sessions:"+token, operator.ID, "EX", 3600).Err
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = s.UpdateSignInTimes(operator.ID)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//更新登录时间
func (s *OperatorService) UpdateSignInTimes(operatorID string) error {
	key := app.System + ":op:" + operatorID + ":sign-in-times"
	err := app.Redis.Cmd("lpush", key, time.Now().Format("2006-01-02 15:04:05")).Err
	if err != nil {
		return errors.WithStack(err)
	}
	err = app.Redis.Cmd("ltrim", key, 0, 2).Err
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

/*//查询用户总数
func (s *OperatorService) CountUsers(cond form.QueryUserRequest) (n int,err error) {
	sess := app.DB.Select("count(*)").From(entity.TableUser)
	if cond.ID != "" {
		sess.AndWhere(dbx.HashExp{"id": cond.ID})
	}
	if string(cond.Gender) != "" {
		sess.AndWhere(dbx.HashExp{"gender": cond.Gender})
	}
	if string(cond.State) != "" {
		sess.AndWhere(dbx.HashExp{"state": cond.State})
	}
	if cond.Oder != "" {
		if cond.Oder == "1"{
			sess.OrderBy("create_time desc")	//降序
		}else{
			sess.OrderBy("create_time asc")		//升序
		}
	}
	//记录总数
	err = sess.Row(&n)
	if err != nil {
		err = errors.Wrap(err, "fail to query devices.")
		return
	}
	return
}*/

//查询用户列表
func (s *OperatorService) QueryBlogUser(cond form.QueryUserRequest) (res []entity.User,err error) {
	sess := app.DB.Select("*").From(entity.TableUser)
	if cond.ID != "" {
		sess.AndWhere(dbx.HashExp{"id": cond.ID})
	}
	if string(cond.Gender) != "" {
		sess.AndWhere(dbx.HashExp{"gender": cond.Gender})
	}
	if string(cond.State) != "" {
		sess.AndWhere(dbx.HashExp{"state": cond.State})
	}
	if cond.Oder != "" {
		if cond.Oder == "1"{
			sess.OrderBy("create_time desc")	//降序
		}else{
			sess.OrderBy("create_time asc")		//升序
		}
	}
	err = sess.All(&res)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if res == nil {
		res = make([]entity.User, 0)
	}
	return
}

//更改用户状态
func (s *OperatorService) ModifyUserStatus(req entity.ModifyUserStateRequest) (user entity.User,err error) {
	err = v.ValidateStruct(&req,
		v.Field(&req.ID, v.Required),
		v.Field(&req.State, v.Required,v.In(entity.UserStateBaned,entity.UserStateOK)),
	)
	if err != nil {
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
	user.State = req.State
	err = app.DB.Model(&user).Update("State")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//获取登录时间
func (s *OperatorService) GetSignInTimes(operatorID string) (times []string, err error) {
	key := app.System + ":op:" + operatorID + ":sign-in-times"
	times, err = app.Redis.Cmd("lrange", key, 0, 2).List()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//验证Token
func (s *OperatorService) CheckToken(token string) (operator entity.Operator, err error) {
	key := "go-blog:op:sessions:" + token
	ID, err := app.Redis.Cmd("GET", key).Str()
	if err != nil {
		if err == redis.ErrRespNil {
			app.Logger.Info().Msg("token expired.")
			err = code.New(http.StatusBadRequest,code.CodeTokenNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}

	err = app.DB.Select().Where(dbx.HashExp{"id":ID}).One(&operator)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusNotFound, code.CodeUserNotExist)
			return
		}
		err = errors.Wrap(err, "fail to find user")
		return
	}

	if operator.State != entity.OperatorStateEnabled {
		app.Logger.Info().Msg("operator status is no enabled.")
		err = code.New(http.StatusBadRequest,code.CodeStateInvalid)
		return
	}

	err = app.Redis.Cmd("EXPIRE", key, 3600).Err
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}
