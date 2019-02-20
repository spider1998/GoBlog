package entity

type UserState int

const TableUser = "user"

const (
	UserStateOK    UserState = 1
	UserStateBaned           = 2
)

type UserGender int8

const(
	UserGenderMale		UserGender	= 1+iota
	UserGenderFemale
)

func (us UserState) Text() string {
	switch us {
	case UserStateOK:
		return "未注销"
	case UserStateBaned:
		return "已注销"
	default:
		return "-"
	}
}

type User struct {
	BaseUser               //基本信息
	State        UserState `json:"state"`         //账号状态（是否已注销）
	PasswordHash []byte    `json:"password_hash"` //密码
	PersonInfo             //个人资料
	AccountInfo            //账户资料
	Mobile       string    `json:"mobile"` //电话
	DatetimeAware
}

type PersonInfo struct {
	Email    string `json:"email"`     //邮箱
	Gender   int    `json:"gender"`    //性别（1:男 2:女）
	Birthday string `json:"birthday"`  //生日
	RealName string `json:"real_name"` //真实姓名
	Area     string `json:"area"`      //地区
}

type AccountInfo struct {
	HeadImg string 		`json:"head_img"` //头像
	Motto   string 		`json:"motto"`    //个性签名
}

type RegisterUserRequest struct {
	Name     string `json:"name"`     //昵称
	Password string `json:"password"` //密码
	Email    string `json:"email"`    //邮箱
	Cach     string `json:"cach"`     //验证码
}

type LoginUserRequest struct {
	Name     string `json:"name"`     //昵称
	Password string `json:"password"` //密码
}

type InfoUpdateRequest struct {
	ID string `json:"id"`
	AccountInfo
	PersonInfo
}

type SetUserPassRequest struct {
	ID          string `json:"id"`
	Password    string `json:"password"`     //密码
	NewPassword string `json:"new_password"` //新密码

}

type BaseUser struct {
	ID  	string 		`json:"id"`                          					//ID
	Name 	string 		`json:"name" gorm:"not null;unique"` 					//昵称
	Tag 	string 		`json:"tag" gorm:"not null;index"`	  			//标签
}

type Contact struct {
	UserID 	string 	`json:"user_id"`
	Name 	string 	`json:"name"`
	Email 	string 	`json:"email"`
	Mobile 	string 	`json:"mobile"`
	Message string 	`json:"message"`
}

type QueryBlogUserResponse struct {
	User 	[]User
	Count 	int `json:"count"`
}

type ModifyUserStateRequest struct {
	ID 		string 		`json:"id"`
	State 	UserState 	`json:"state"`
}


func (User) TableName() string {
	return TableUser
}
