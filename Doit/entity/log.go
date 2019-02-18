package entity

import "time"

type LogUserType int

const (
	LogUserTypeUser     LogUserType = 1
	LogUserTypeOperator LogUserType = 2
)

type Log struct {
	ID         string                 `json:"id" xorm:"pk"`
	UserType   LogUserType            `json:"user_type"`
	UserID     string                 `json:"user_id"`
	UserName   string                 `json:"user_name"`
	System     string                 `json:"system"`
	Action     string                 `json:"action"`
	Remark     string                 `json:"remark"`
	IP         string                 `json:"ip"`
	Ext        map[string]interface{} `json:"ext"`
	CreateTime time.Time              `json:"create_time" xorm:"created"`
	UpdateTime time.Time              `json:"update_time" xorm:"updated"`
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
