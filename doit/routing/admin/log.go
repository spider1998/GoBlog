package admin

import (
	"strconv"
	"github.com/go-ozzo/ozzo-routing"
	"Project/doit/form"
	"Project/doit/code"
	"net/http"
	"Project/doit/entity"
	"Project/doit/util"
	"Project/doit/service"
)

type LogHandler struct{}

//查询日志
func (LogHandler) QueryLogs(c *routing.Context) error {
	var cond form.QueryLogsCond
	if userTypeStr := c.Query("user_type"); userTypeStr != "" {
		userType, err := strconv.Atoi(userTypeStr)
		if err != nil {
			return code.New(http.StatusBadRequest,code.CodeInvalidData)
		}
		cond.UserType = entity.LogUserType(userType)
	}
	cond.Remark = c.Query("remark")
	cond.FromTime = c.Query("from_time")
	cond.ToTime = c.Query("to_time")

	n, err := service.Log.CountLogs(cond)
	if err != nil {
		return err
	}
	pages := util.NewPagerFromRequest(c, n)
	logs, err := service.Log.QueryLogs(pages.Offset(), pages.Limit(), cond)
	if err != nil {
		return err
	}

	return c.Write(logs)
}
