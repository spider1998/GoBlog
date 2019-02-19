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
	Ext        map[string]interface{} `json:"ext" gorm:"type:json"`
	CreateTime time.Time              `json:"create_time" gorm:"created"`
	UpdateTime time.Time              `json:"update_time" gorm:"updated"`
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

func (l *Log) BeforeInsert() {
	if l.Ext == nil {
		l.Ext = make(map[string]interface{})
	}
}

func (l *Log) BeforeUpdate() {
	if l.Ext == nil {
		l.Ext = make(map[string]interface{})
	}
}

func (Log) TableName() string {
	return TableLog
}
