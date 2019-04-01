package entity

const TableMessage = "message"

type MessageType int

type MessageAuthType int8

const (
	MessageTypeAuth MessageType = 1 + iota
	MessageTypeNotice
	MessageAuthTypeAllowed MessageAuthType = 1 + iota
	MessageAuthTypeRefused
)

type Message struct {
	ID       string      `json:"id"`         //消息记录ID
	UserID   string      `json:"user_id"` //用户ID
	OriginID string      `json:"origin_id"`  //来源用户ID
	Type     MessageType `json:"type"`       //消息类型
	ServerID string      `json:"server_id"`  //服务ID（授权，申请）
	Title    string      `json:"title"`      //消息标题
	Content  string      `json:"content"`    //消息内容
	Read     bool        `json:"read"`       //阅读状态
	DatetimeAware
}

func (Message) TableName() string {
	return TableMessage
}
