package form

type AddFriendRequest struct {
	UserID   string `json:"user_id"`   //用户ID
	Name 	string `json:"name"`	//姓名
	FriendID string `json:"friend_id"` //好友ID
	Reason   string `json:"reason"`    //申请理由
}
