package friend

import (
	"Project/doit/app"
	"Project/doit/code"
	"Project/doit/entity"
	"Project/doit/form"
	"Project/doit/handler/session"
	"Project/doit/service"
	"Project/doit/util"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/pkg/errors"
	"net/http"
)

//拉黑好友
func PullBlack(c *routing.Context) (err error) {
	userID := session.GetUserSession(c).ID
	recID := c.Query("record_id")
	state := c.Query("state")
	err = service.Friend.PullBlack(userID, recID,state)
	if err != nil {
		return
	}
	return
}

//获取好友列表
func GetFriendList(c *routing.Context) (err error) {
	userID := c.Query("user_id")
	state := c.Query("state")
	friends, err := service.Friend.GetFriendList(userID, state)
	if err != nil {
		return err
	}
	var res []entity.Friend
	pager := util.GetPaginatedListFromRequest(c, len(friends))
	if pager.Offset()+pager.Limit() <= pager.TotalCount {
		res = friends[pager.Offset() : pager.Offset()+pager.Limit()]
	} else {
		res = friends[pager.Offset():pager.TotalCount]
	}
	return c.Write(res)

}

//删除好友
func DeleteFriend(c *routing.Context) (err error) {
	userID := session.GetUserSession(c).ID
	recID := c.Query("record_id")
	var rec entity.Friend
	err = app.DB.Select().Where(dbx.HashExp{"id": recID}).One(&rec)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeRecordNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	if userID != rec.UserID && userID != rec.FriendID {
		err = code.New(http.StatusBadRequest, code.CodeUserAccessSessionInvalid)
		return
	}
	err = app.DB.Delete("friend", dbx.NewExp("id={:id}", dbx.Params{"id": recID})).
		All(&rec)
	if err != nil {
		return err
	}
	return
}

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
