package service

import (
	"Project/doit/app"
	"Project/doit/code"
	"Project/doit/entity"
	"Project/doit/form"
	"Project/doit/util"
	"github.com/go-ozzo/ozzo-dbx"
	v "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"net/http"
	"strconv"
	"time"
)

var Friend = FriendService{}

type FriendService struct{}

//好友申请授权
func (f *FriendService) AddAuthorization(req form.AddFriendRequest, state string) (err error) {
	err = v.ValidateStruct(&req,
		v.Field(&req.Name, v.Required),
		v.Field(&req.UserID, v.Required),
		v.Field(&req.FriendID, v.Required),
		v.Field(&req.Reason, v.Required),
	)
	if err != nil {
		return
	}
	status, err := strconv.Atoi(state)
	if err != nil {
		return
	}
	if status == int(entity.AcceptApp) {
		var friend entity.Friend
		friend.ID = xid.New().String()
		friend.FriendID = req.FriendID
		friend.UserID = req.UserID
		friend.UserState = entity.FriendOK
		friend.FriendState = entity.FriendOK
		friend.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		err = app.DB.Transactional(func(tx *dbx.Tx) error {
			err = tx.Model(&friend).Insert()
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			if util.IsDBDuplicatedErr(err) {
				err = code.New(http.StatusConflict, code.CodeFriendExist)
				return
			}
			err = errors.Wrap(err, "fail to create friend relationship")
			return
		}
	}
	return

}

//查询人员
func (f *FriendService) QueryUsers(req entity.QueryUserRequest) (persons []entity.BaseUser, err error) {
	query := app.DB.Select("*").From(entity.TableUser)

	if req.Name != "" {
		query.AndWhere(dbx.Like("name", req.Name))
	}
	if req.Gender == entity.UserGenderMale {
		query.AndWhere(dbx.HashExp{"gender": entity.UserGenderMale})
	}
	if req.Gender == entity.UserGenderFemale {
		query.AndWhere(dbx.HashExp{"gender": entity.UserGenderFemale})
	}
	if req.Tag != "" {
		query.AndWhere(dbx.HashExp{"tag": req.Tag})
	}
	var users []entity.User
	var person entity.BaseUser
	err = query.All(&users)
	if err != nil {
		err = errors.Wrap(err, "fail to query users.")
		return
	}
	for _, user := range users {
		person.ID = user.ID
		person.Tag = user.Tag
		person.Name = user.Name
		persons = append(persons, person)
	}
	return
}
