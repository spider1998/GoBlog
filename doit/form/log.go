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
	UserType entity.LogUserType
	Remark   string
	FromTime string
	ToTime   string
}

type QueryUserRequest struct {
	ID		string `json:"id"`
	State 	entity.UserState `json:"state"`
	Gender	entity.UserGender `json:"gender"`
	Oder 	string `json:"oder"`
}

type QueryArticleRequest struct {
	ID 		string `json:"id"`
	Sort 	string `json:"sort"`
}

type QueryArticleResponse struct {
	ID 		string `json:"id"`
	Sort 	string `json:"sort"`
	Auth 	string `json:"auth"`
	Title  	string `json:"title"`
	entity.DatetimeAware
}

type GetArticlesResponse struct {
	Count 	int `json:"count"`
	Arts 	[]QueryArticleResponse `json:"arts"`
}



