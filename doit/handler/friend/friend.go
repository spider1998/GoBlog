package friend

import (
	"Project/doit/app"
	"Project/doit/code"
	"Project/doit/entity"
	"Project/doit/form"
	"Project/doit/handler/session"
	"Project/doit/service"
	"Project/doit/util"
	"github.com/go-ozzo/ozzo-routing"
	"net/http"
)

//好友申请授权
func AddAuthorization(c *routing.Context) (err error) {
	var req form.AddFriendRequest
	err = c.Read(&req)
	if err != nil {
		return
	}
	if session.GetUserSession(c).ID != req.FriendID {
		return code.New(http.StatusBadRequest, code.CodeUserAccessSessionInvalid)
	}
	state := c.Query("state")
	err = service.Friend.AddAuthorization(req, state)
	if err != nil {
		return
	}
	var content string
	if state == "1" {
		content = "申请通过！好友【" + session.GetUserSession(c).Name + "】添加成功！"
	} else {
		content = "申请拒绝！好友【" + session.GetUserSession(c).Name + "】拒绝好友申请！"
	}
	err = service.Message.Create(req.UserID, app.Conf.FriendNotice, req.FriendID, content, req.FriendID)
	if err != nil {
		return
	}
	return
}

//添加好友申请
func AddFriends(c *routing.Context) (err error) {
	var req form.AddFriendRequest
	err = c.Read(&req)
	if err != nil {
		return
	}
	err = service.Message.Create(req.FriendID, app.Conf.FriendNotice, req.UserID,
		req.Name+" 请求添加好友： 申请理由："+req.Reason, req.UserID)
	if err != nil {
		return
	}
	return
}

//查询人员
func QueryUsers(c *routing.Context) (err error) {
	var req entity.QueryUserRequest
	req.Name = c.Param("name")
	gender := c.Param("gender")
	req.Tag = c.Param("tags")
	if gender == "1" {
		req.Gender = entity.UserGenderMale
	} else {
		req.Gender = entity.UserGenderFemale
	}
	persons, err := service.Friend.QueryUsers(req)
	if err != nil {
		return
	}
	if len(persons) == 0 {
		persons = []entity.BaseUser{}
	}
	pager := util.GetPaginatedListFromRequest(c, len(persons))
	if pager.Offset()+pager.Limit() <= pager.TotalCount {
		return c.Write(persons[pager.Offset() : pager.Offset()+pager.Limit()])
	} else {
		return c.Write(persons[pager.Offset():pager.TotalCount])
	}
}
