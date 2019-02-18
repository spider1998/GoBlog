package admin

import (
	"fmt"
	"github.com/go-ozzo/ozzo-routing"
	"Project/Doit/form"
	"Project/Doit/code"
	"net/http"
	"Project/Doit/service"
	"Project/Doit/app"
	"github.com/go-ozzo/ozzo-routing/access"
	"Project/Doit/util"
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
