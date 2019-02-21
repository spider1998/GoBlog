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


