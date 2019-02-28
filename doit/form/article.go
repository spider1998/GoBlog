package form

import "Project/doit/entity"

type CommentArticleRequest struct {
	ArtID	string `json:"art_id"`
	UserID 	string `json:"user_id"`
	Name	string `json:"name"`
	Content string `json:"content"`
}

type CommentReplyRequest struct {
	ComID		string 		`json:"com_id"`
	FatherID	string 		`json:"father_id"`
	FatherName	string 		`json:"father_name"`
	UserID		string 		`json:"user_id"`
	Name		string 		`json:"name"`
	Content 	string 		`json:"content"`
}

type ArticleCommentResponse struct {
	Comment MainComment 	`json:"comment"`
	Replies []SonReply 		`json:"replies"`
}

type SortStaticResponse struct {
	Sorts 	[]string 	`json:"sorts"`
	Arry	[]int 		`json:"arry"`
}

type GenderStaticResponse struct {
	Male	[]int `json:"male"`
	Female	[]int `json:"female"`
}

type AreaStatisticResponse struct {
	Area 	[]string 	`json:"area"`
	Array 	[]int 		`json:"array"`
}

type ArticleTopResponse struct {
	Hot 	[10]HotTop  	`json:"hot"`
	Read	[10]ReadTop 	`json:"read"`
}

type HotTop struct {
	ID		string 		`json:"id"`
	Title 	string 		`json:"title"`
	Auth 	string 		`json:"auth"`
	Hot 	int 		`json:"hot"`
	Time    string 		`json:"time"`
}

type ReadTop struct {
	ID		string 		`json:"id"`
	Title 	string 		`json:"title"`
	Auth 	string 		`json:"auth"`
	Read 	int 		`json:"read"`
	Time    string 		`json:"time"`
} 



type MainComment struct {
	ComID 			string `json:"com_id"`
	UserID 			string `json:"user_id"`
	Name 			string `json:"name"`
	Content 		string `json:"content"`
	ReplyCount		int `json:"reply_count"`
	entity.DatetimeAware
}

type SonReply struct {
	entity.ReplyBase
}


