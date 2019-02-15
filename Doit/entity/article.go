package entity


const TableArticle = "article"
const TableArticleVersion = "article_version"
const TableArticleForward = "article_forward"

type ModifyType int

const(
	ModifyTypeAble  ModifyType		=	1+iota
	ModifyTypeEnable
)
type Article struct {
	ID          	string     	`json:"id" gorm:"index;not null"`           //唯一ids
	UserId      	string     	`json:"user_id"`      						//用户id
	PartPersons 	string 		`json:"part_persons"` 						//贡献者
	Version			int			`json:"version"`							//文章版本
	ModifyType		ModifyType 		`json:"modify_type"`					//文章修改类型
	BaseArticle            													//文章基本字段
	ArticleContent
	Comment
	DatetimeAware
}

type ArticleVersion struct {
	ID          string     `json:"id" gorm:"index;not null"` //唯一ids
	ArtID			string `json:"art_id"`
	UserId      	string     	`json:"user_id"`      						//用户id
	Version			int			`json:"version"`							//文章版本
	ModifyType		ModifyType 		`json:"modify_type"`					//文章修改类型
	BaseArticle            													//文章基本字段
	ArticleContent
	Comment
	DatetimeAware

}

type BaseArticle struct {
	Title   string `json:"title"`                          //标题
	Auth    string `json:"auth"`                           //主作者
	Sort    string `json:"sort"`                           //类别
	Content string `json:"content" gorm:"type:text;not null;"` //最新内容
}

type ArticleContent struct {
	SecondTitle string      `json:"second_title"` //副标题
	Photo       string 		`json:"photo"`        //图片
	Attachment  string      `json:"attachment"`   //附件
	Hot         int      	`json:"hot"`          //热度
	Forward     string      `json:"forward"`      //转发数
	Read 		int 		`json:"read"`		  //阅读量
}

type Comment struct {
	Commentator   string `json:"commentator"` //评论员
	Comments      string `json:"comments"`    //评论内容
	ComUpdateTime string `json:"update_time"` //评论时间
}

type BasePhoto struct {
	Url string `json:"url"` //图片链接
}

type CreateArticleRequest struct {
	UserId string `json:"user_id"` 				//用户ID
	ModifyType		ModifyType `json:"modify_type"`			//文章令牌
	BaseArticle
	Photo       string 		`json:"photo"`        //图片
	SecondTitle string      `json:"second_title"` //副标题
	Attachment  string      `json:"attachment"`   //附件

}

type VerifyArticleRequest struct { //
	ID     string `json:"id"`      //文章ID
	UserId string `json:"user_id"` 				//用户ID
	ModifyType		ModifyType `json:"modify_type"`			//文章令牌
	BaseArticle
	Photo       string 		`json:"photo"`        //图片
	SecondTitle string      `json:"second_title"` //副标题
	Attachment  string      `json:"attachment"`   //附件
}

type UpdateArticleRequest struct {
	ID      string `json:"id"`                             //文章ID
	Content string `json:"content" gorm:"not null;unique"` //内容
}

type RestoreArticleRequest struct {
	ArtId		string 		`json:"art_id"`		 //	文章ID
	UserId      string     `json:"user_id"`      //用户id
	Version		int		`json:"version"`		//文章版本
}

type ArticleResponse struct {
	ID          string     `json:"id" gorm:"index;not null"` //唯一ids
	UserId      string     `json:"user_id"`                  //用户id
	PartPersons string     `json:"part_persons"`             //贡献者
	Version     int        `json:"version"`                  //文章版本
	ModifyType  ModifyType `json:"modify_type"`              //文章修改类型
	SecondTitle string      `json:"second_title"` //副标题
	Photo       []byte		`json:"photo"`        //图片
	Attachment  string      `json:"attachment"`   //附件
	Hot         int      	`json:"hot"`          //热度
	Forward     string      `json:"forward"`      //转发数
	BaseArticle                                              //文章基本字段
	Comment
	DatetimeAware
}

type ArticleForwardRequest struct {
	ArtID string `json:"art_id"`
	UserID 	string `json:"user_id"`
	Reason	string `json:"reason"`
}

type StateForward int8

const(
	StateForwardWaite	StateForward	=	1 + iota
	StateForwardRefused
	StateForwardFinished
)

type ArticleForward struct {
	ID			string 				`json:"id"`
	ArtID 		string 				`json:"art_id" gorm:"index;not null"`
	AuthID 		string 				`json:"auth_id"`
	ForwardID	string 				`json:"forward_id"`
	Reason		string 				`json:"reason" gorm:"type:text;not null"`
	Status 		StateForward 		`json:"status"`
	DatetimeAware
}

type ArticleAuthorazation struct {
	ArtID	string `json:"art_id"`
	RecordID string `json:"record_id"`
	State 	StateForward `json:"state"`
}






func (Article) TableName() string {
	return TableArticle
}
func (ArticleVersion) TableName() string {
	return TableArticleVersion
}
func (ArticleForward) TableName() string {
	return TableArticleForward
}
