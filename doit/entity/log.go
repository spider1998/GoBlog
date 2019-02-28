package entity

import "time"

const TableLog = "log"

type LogUserType int

const (
	LogUserTypeUser     LogUserType = 1
	LogUserTypeOperator LogUserType = 2
)

type Log struct {
	ID         string                 `json:"id" gorm:"index"`
	UserType   LogUserType            `json:"user_type"`
	UserID     string                 `json:"user_id"`
	UserName   string                 `json:"user_name"`
	System     string                 `json:"system"`
	Action     string                 `json:"action"`
	Remark     string                 `json:"remark"`
	IP         string                 `json:"ip"`
	CreateTime time.Time              `json:"create_time" gorm:"created"`
}

type CreateLogRequest struct {
	Token    string                 `json:"token"`
	UserType LogUserType            `json:"user_type"`
	System   string                 `json:"system"`
	Action   string                 `json:"action"`
	IP       string                 `json:"ip"`
	Remark   string                 `json:"remark"`
	Ext      map[string]interface{} `json:"ext"`
}


func (Log) TableName() string {
	return TableLog
}
