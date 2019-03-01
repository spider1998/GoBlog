package form

import (
	"Project/doit/entity"
)


const (
	LogUserTypeUser     entity.LogUserType = 1
	LogUserTypeOperator entity.LogUserType = 2
)

type CreateLogRequest struct {
	Token    string                 `json:"token"`
	UserType entity.LogUserType     `json:"user_type"`
	System   string                 `json:"system"`
	Action   string                 `json:"action"`
	IP       string                 `json:"ip"`
	Remark   string                 `json:"remark"`
	Ext      map[string]interface{} `json:"ext"`
}

type QueryLogsCond struct {
	UserType entity.LogUserType `json:"user_type"`
	Remark   string `json:"remark"`
	FromTime string `json:"from_time"`
	ToTime   string `json:"to_time"`
}



