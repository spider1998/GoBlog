package admin

import (
	"Project/doit/code"
	"Project/doit/entity"
	"Project/doit/handler/article"
	"Project/doit/service"
	"github.com/go-ozzo/ozzo-routing"
	"net/http"
	"strings"
)

func ManagerRegisterRoutes(router *routing.RouteGroup) {
	var (
		operatorHandler = new(OperatorHandler)
		captchaHandler  = new(CaptchaHandler)
		logHandler      = new(LogHandler)
	)

	{
		router.Get("/captcha", captchaHandler.Generate)  // 拉取图片验证码
		router.Post("/sessions", operatorHandler.SignIn) // 管理员登录
		router.Get("/announcements", operatorHandler.GetAnnouncements)    //获取公告
	}

	router.Use(sessionChecker)

	{
		/*-----------------------------------------System------------------------------------------------*/
		router.Post("/announcements", operatorHandler.CreateSiteAnnounce) // 发布公告
		router.Get("/backups", listBackups)                               // 获取备份列表
		router.Post("/backups", makeBackup)                               // 备份数据库
		router.Post("/restores", restoreBackup)                           // 还原数据库
		router.Patch("/schedules/<key>", updateSchedule)                  // 更新计划任务
		router.Get("/schedules/<key>", getSchedule)                       // 获取计划任务
		/*-----------------------------------------Statistics------------------------------------------------*/
		router.Get("/statistics/month/<year>", operatorHandler.GetMonthArticle) // 获取每个月份文章发布数
		router.Get("/statistics", operatorHandler.GetStatistics)                // 获取站点统计数据
		router.Get("/statistics/sort", operatorHandler.GetSortStatistic)        // 获取文章各类别统计
		router.Get("/statistics/gender", operatorHandler.GetGenderStatic)       // 获取性别各时间段发文统计
		router.Get("/statistics/area", operatorHandler.GetAreaStatic)           // 获取用户地区分布统计
		router.Get("/articles/top10", operatorHandler.GetArticleTop)            // 获取文章排行前十
		/*-----------------------------------------Log------------------------------------------------*/
		router.Get("/logs", logHandler.QueryLogs) // 查询日志
		/*-----------------------------------------User------------------------------------------------*/
		router.Get("/sessions/current", operatorHandler.GetSession)     // 获取管理员信息
		router.Get("/users", operatorHandler.QueryBlogUser)             // 查询用户
		router.Patch("/users/status", operatorHandler.ModifyUserStatus) // 启用/禁用用户
		/*-----------------------------------------Article------------------------------------------------*/
		router.Get("/articles/<art_id>", article.GetArticle)                       //	获取指定文章
		router.Get("/articles", operatorHandler.GetArticlesList)                   //	获取文章列表
		router.Delete("/articles/<art_id>", operatorHandler.DeleteArticle)         //	删除指定文章
		router.Post("/sorts", operatorHandler.CreateArticleSort)                   //	创建文章分类
		router.Patch("/sorts/status/<sort_id>", operatorHandler.ModifyArticleSort) //	删除文章分类
		router.Get("/sorts", operatorHandler.GetArticlesSorts)                     //	获取文章分类

	}
}

const (
	sessionTokenHeaderKey = "X-Access-Token"
	sessionKey            = "session.operator"
)

func sessionChecker(c *routing.Context) error {
	token := c.Request.Header.Get(sessionTokenHeaderKey)
	if token == "" {
		token = c.Query(strings.ToLower(sessionTokenHeaderKey))
		if token == "" {
			return code.New(http.StatusNotFound, code.CodeOperatorTokenRequired)
		}
	}
	operator, err := service.Operator.CheckToken(token)
	if err != nil {
		return err
	}
	c.Set(sessionKey, operator)
	return c.Next()
}

func getSessionOperator(c *routing.Context) entity.Operator {
	return c.Get(sessionKey).(entity.Operator)
}
