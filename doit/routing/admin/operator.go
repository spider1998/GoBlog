package admin

import (
	"fmt"
	"github.com/go-ozzo/ozzo-routing"
	"Project/doit/form"
	"Project/doit/code"
	"net/http"
	"Project/doit/service"
	"Project/doit/app"
	"github.com/go-ozzo/ozzo-routing/access"
	"Project/doit/util"
	"Project/doit/entity"
)

type OperatorHandler struct{}

func (OperatorHandler) SignIn(c *routing.Context) error {
	var request form.OperatorSignInRequest
	err := c.Read(&request)
	if err != nil {
		return code.New(http.StatusBadRequest,code.CodeInvalidData)
	}
	token, operator, err := service.Operator.SignIn(request)
	if err != nil {
		return err
	}

	service.Log.LogOperator(
		operator,
		app.System,
		"operator.sign-in",
		fmt.Sprintf("管理员登录。"),
		access.GetClientIP(c.Request),
		util.M{"operator": operator},
	)

	return c.Write(map[string]string{"token": token})
}

func (OperatorHandler) GetSession(c *routing.Context) error {
	operator := getSessionOperator(c)
	session := entity.OperatorSession{
		Operator: operator,
	}
	times, err := service.Operator.GetSignInTimes(operator.ID)
	if err != nil {
		return err
	}
	if len(times) > 0 {
		session.SignInTime = times[0]
	}
	if len(times) > 1 {
		session.LastSignInTime = times[1]
	}
	return c.Write(session)
}

//查询用户
func (OperatorHandler) QueryBlogUser(c *routing.Context) error {
	var req form.QueryUserRequest
	req.ID = c.Query("user_id")
	req.Oder = c.Query("oder")
	if c.Query("gender") == "1"{
		req.Gender = 1
	}else if c.Query("gender") == "2"{
		req.Gender = 2
	}
	if c.Query("state") == "1"{
		req.State = 1
	}else if c.Query("state") == "2"{
		req.State = 2
	}
	response,err := service.Operator.QueryBlogUser(req)
	if err != nil{
		return err
	}
	if len(response) == 0{
		response = []entity.User{}
	}
	var users []entity.User
	pager := util.GetPaginatedListFromRequest(c, len(response))
	if pager.Offset()+pager.Limit() <= pager.TotalCount {
		users = response[pager.Offset() : pager.Offset()+pager.Limit()]
	} else {
		users = response[pager.Offset():pager.TotalCount]
	}
	var res entity.QueryBlogUserResponse
	res.User = users
	res.Count = len(response)
	return c.Write(res)

}

//修改用户账号状态
func (OperatorHandler) ModifyUserStatus(c *routing.Context) error {
	state := c.Query("state")
	var req entity.ModifyUserStateRequest
	if state == "1"{
		req.State = entity.UserStateOK
	}else if state == "2"{
		req.State = entity.UserStateBaned
	}
	req.ID = c.Query("user_id")
	user,err := service.Operator.ModifyUserStatus(req)
	if err != nil{
		return err
	}
	service.Log.LogOperator(
		getSessionOperator(c),
		app.System,
		"operator.modify-user-status",
		fmt.Sprintf("修改用户账号状态。"),
		access.GetClientIP(c.Request),
		util.M{"operator": user},
	)
	return c.Write(http.StatusOK)
}

//获取文章列表（条件查询）
func (OperatorHandler) GetArticlesList(c *routing.Context) error {
	var req form.QueryArticleRequest
	req.ID = c.Query("art_id")
	req.Sort = c.Query("sort")
	articles,err := service.Operator.GetArticlesList(req)
	if err != nil{
		return err
	}
	var arts []form.QueryArticleResponse
	pager := util.GetPaginatedListFromRequest(c, len(articles))
	if pager.Offset()+pager.Limit() <= pager.TotalCount {
		arts = articles[pager.Offset() : pager.Offset()+pager.Limit()]
	} else {
		arts = articles[pager.Offset():pager.TotalCount]
	}

	var res form.GetArticlesResponse
	res.Arts = arts
	res.Count = len(articles)
	return c.Write(res)
}

