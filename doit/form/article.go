package form

type CommentArticleRequest struct {
	ArtID	string `json:"art_id"`
	UserID 	string `json:"user_id"`
	Name	string `json:"name"`
	Content string `json:"content"`
}

type CommentReplyRequest struct {
	ComID		string 		`json:"com_id"`
	UserID		string 		`json:"user_id"`
	Name		string 		`json:"name"`
	Content 	string 		`json:"content"`
}

