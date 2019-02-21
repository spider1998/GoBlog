package entity

const TableComment  = "comment"
const TableReply  = "reply"

type Comment struct {
	ID 			string 		`json:"id"`				//记录ID
	ArtID		string 		`json:"art_id"`			//文章ID
	UserID		string 		`json:"user_id"`		//评论用户ID
	Name		string 		`json:"name"`			//评论用户姓名
	ReplyCount	int 		`json:"reply_count"`	//回复数
	Content 	string 		`json:"content"`		//评论内容
	DatetimeAware									//时间

}

type Reply struct {
	ID 			string 		`json:"id"`				//记录ID
	ComID 		string 		`json:"com_id"`			//评论ID
	ReplyBase
}

type ReplyBase struct {
	FatherID	string 		`json:"father_id"`		//父评论用户ID
	FatherName	string 		`json:"father_name"`	//父评论用户姓名
	UserID 		string 		`json:"user_id"`		//用户ID
	Name 		string 		`json:"name"`			//用户名
	Content 	string 		`json:"content"`		//回复内容
	DatetimeAware
}


func (Comment) TableName() string {
	return TableComment
}

func (Reply) TableName() string {
	return TableReply
}
