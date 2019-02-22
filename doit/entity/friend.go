package entity

const TableFriend = "friend"

type FriendStatus int

const (
	FriendOK    FriendStatus = 1
	FriendBlack              = 2
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
	ID     string       `json:"id"`     //ID
	Name   string       `json:"name"`   //姓名
	Uid    int          `json:"uid"`    //所属好友
	Status FriendStatus `json:"status"` //是否拉黑
}

type QueryUserRequest struct {
	Name   string     `json:"name"`
	Gender UserGender `json:"gender"`
	Tag    string     `json:"tags"`
}

func (Friend) TableName() string {
	return TableFriend
}
