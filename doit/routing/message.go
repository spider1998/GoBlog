package routing

import (
	"net/http"

	"strconv"

	"github.com/go-ozzo/ozzo-dbx"
	"github.com/pkg/errors"

	"Project/doit/app"
	"Project/doit/code"
	"Project/doit/entity"
	"Project/doit/handler/session"
	"Project/doit/service"
	"Project/doit/util"
	"github.com/go-ozzo/ozzo-routing"
)

//获取消息列表
func getMessagesList(c *routing.Context) error {
	//未读、已读
	readStatus := c.Query("read_status")
	var read bool
	if readStatus == "read" {
		read = true
	}
	//查询记录总数
	query := app.DB.Select("count(*)").From(entity.TableMessage).Where(dbx.HashExp{
		"user_id": session.GetUserSession(c).ID,
		"read":    read,
	})

	var cnt int
	err := query.Row(&cnt)
	if err != nil {
		return errors.Wrap(err, "fail to query message.")
	}
	c.Response.Header().Set("X-Total-Count", strconv.Itoa(cnt))

	//解析页码和每页显示记录数
	page, size, err := util.ParsePagination(c)
	if err != nil {
		return errors.Wrap(err, "fail to parse pagination.")
	}

	var msg []entity.Message
	err = query.Select().Limit(size).Offset((page - 1) * size).OrderBy("create_time desc").All(&msg)
	if err != nil {
		return errors.Wrap(err, "fail to select message.")
	}

	return c.Write(msg)
}

//修改消息状态
func changeMessageStatus(c *routing.Context) error {
	msgId := c.Query("message_id")
	if msgId == "" {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err("request content is empty")
	}
	//默认为true，后续可以更改
	session := session.GetUserSession(c)
	message, err := service.Message.Read(session.ID, msgId, true)
	if err != nil {
		return err
	}
	return c.Write(message)
}
