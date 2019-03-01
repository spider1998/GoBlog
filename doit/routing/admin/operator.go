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
	"strconv"
	"time"
)

type OperatorHandler struct{}

//获取站点统计数据
func (OperatorHandler) GetStatistics(c *routing.Context) error {
	res,err := service.Operator.GetStatistics()
	if err != nil{
		return err
	}
	return c.Write(res)
}

//获取文章类别统计
func (OperatorHandler) GetSortStatistic(c *routing.Context) error {
	res,err := service.Operator.GetSortStatistic()
	if err != nil{
		return err
	}
	return c.Write(res)
}

//获取性别各时间段发文统计
func (OperatorHandler) GetGenderStatic(c *routing.Context) error {
	res,err := service.Operator.GetGenderStatic()
	if err != nil{
		return err
	}
	return c.Write(res)
}

//获取文章排行
func (OperatorHandler) GetArticleTop(c *routing.Context) error {
	res,err := service.Operator.GetArticleTop()
	if err != nil{
		return err
	}
	return c.Write(res)
}

//获取用户地区分布统计
func (OperatorHandler) GetAreaStatic(c *routing.Context) error {
	res,err := service.Operator.GetAreaStatic()
	if err != nil{
		return err
	}
	return c.Write(res)
}

//获取每个月份文章发布数
func (OperatorHandler) GetMonthArticle(c *routing.Context) error {
	yearStr := c.Param("year")
	var year int
	if yearStr == "0"{
		year = time.Now().Year()
	}else{
		yearInt,err := strconv.Atoi(yearStr)
		if err != nil{
			return err
		}
		year = yearInt
	}
	res,err := service.Operator.GetMonthArticle(year)
	if err != nil{
		return err
	}
	return c.Write(res)
}

//管理员登录
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
	fmt.Println(response)
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
		fmt.Sprintf("修改用户账号状态id("+user.ID+")。"),
		access.GetClientIP(c.Request),
		util.M{"operator": user},
	)
	return c.Write(http.StatusOK)
}

//删除文章
func (OperatorHandler)DeleteArticle(c *routing.Context) error {
	articleID := c.Param("art_id")
	err := service.Operator.DeleteArticle(articleID)
	if err != nil {
		return err
	}
	_,err = service.Log.LogOperator(
		getSessionOperator(c),
		app.System,
		"operator.create-sort",
		fmt.Sprintf("管理员删除文章id("+articleID+")。"),
		access.GetClientIP(c.Request),
		util.M{"article": articleID},
	)
	if err != nil{
		return err
	}

	return c.Write(http.StatusOK)
}

//获取文章分类
func (OperatorHandler)GetArticlesSorts(c *routing.Context) error {
	sorts,err := service.Operator.GetArticlesSorts()
	if err != nil {
		return err
	}
	return c.Write(sorts)
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

//删除文章分类
func (OperatorHandler)DeleteArticleSort(c *routing.Context) error {
	sortID := c.Param("sort_id")
	sort,err := service.Operator.DeleteArticlesSorts(sortID)
	if err != nil {
		return err
	}
	service.Log.LogOperator(
		getSessionOperator(c),
		app.System,
		"operator.create-sort",
		fmt.Sprintf("管理员删除文章分类。"),
		access.GetClientIP(c.Request),
		util.M{"sort": sort.Name},
	)
	return c.Write(sort)
}


//创建文章分类
func (OperatorHandler)CreateArticleSort(c *routing.Context) error {
	var req form.CreateArticleSortRequest
	req.Name = getSessionOperator(c).Name
	err := c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	sort,err := service.Operator.CreateArticleSort(req)
	if err != nil{
		return err
	}
	service.Log.LogOperator(
		getSessionOperator(c),
		app.System,
		"operator.create-sort",
		fmt.Sprintf("创建新的文章分类。"),
		access.GetClientIP(c.Request),
		util.M{"sort": sort},
	)
	return c.Write(http.StatusOK)
}

