package entity

const TableFriend = "friend"

type FriendStatus int

type AppFriendState int8

const (
	FriendOK FriendStatus = 1 + iota
	FriendBlack
	AcceptApp	AppFriendState	=	1 +iota
	RefusedApp
)

func (fr FriendStatus) Text() string {
	switch fr {
	case FriendOK:
		return "正常"
	case FriendBlack:
		return "已拉黑"
	default:
		return "-"
	}
}

type Friend struct {
	ID          string       `json:"id"`                                         //ID
	UserID      string       `json:"user_id" gorm:"index;not null"`              //用户ID
	FriendID    string       `json:"friend_id" gorm:"index;not null"`            //好友ID
	FriendState FriendStatus `json:"friend_state"`                               //是否拉黑
	UserState   FriendStatus `json:"user_state"`                                 //是否被拉黑
	CreateTime  string       `json:"create_time" gorm:"type:datetime;not null;"` //创建好友时间
}

type QueryUserRequest struct {
	Name   string     `json:"name"`
	Gender UserGender `json:"gender"`
	Tag    string     `json:"tags"`
}

func (Friend) TableName() string {
	return TableFriend
}
