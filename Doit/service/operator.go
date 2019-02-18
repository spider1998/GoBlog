package service

import (
	"Project/Doit/form"
	"Project/Doit/entity"
)

var Operator = &OperatorService{}

type OperatorService struct{}

func (s *OperatorService) SignIn(request form.OperatorSignInRequest) (token string, operator entity.Operator, err error) {
	return
}

func (s *OperatorService) CheckToken(token string) (operator entity.Operator, err error) {
	return
}
