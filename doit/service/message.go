package service

import (
	"Project/doit/app"
	"Project/doit/code"
	"Project/doit/entity"
	"Project/doit/util"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

type MessageService struct{}

var Message MessageService

func (m MessageService) Create(accountID, title, content string) error {
	var message entity.Message
	message.ID = uuid.New().String()
	message.UserID = accountID
	message.Title = title
	message.Content = content
	message.Read = false
	message.DatetimeAware = entity.DatetimeAwareNow()
	err := app.DB.Model(&message).Insert()
	if err != nil {
		app.Logger.Error().Err(err).Msg("fail to insert message to db.")
		return err
	}
	return nil
}

//路由校验，提取消息uid,msgId和状态，调用该函数修改状态
func (m MessageService) Read(uid, msgId string, read bool) (message entity.Message, err error) {
	err = app.DB.Select().Where(dbx.HashExp{"id": msgId, "account_id": uid}).One(&message)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusNotFound, code.CodeMessageNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}

	message.Read = read
	message.UpdateTime = util.DateTimeStd()
	err = app.DB.Model(&message).Update("Read", "UpdateTime")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
