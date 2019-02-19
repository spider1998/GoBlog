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

