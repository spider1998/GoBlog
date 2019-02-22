package service

import (
	"Project/doit/app"
	"Project/doit/entity"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/pkg/errors"
)

var Friend = FriendService{}

type FriendService struct{}

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
