package entity


const TableArticle = "article"

type Article struct {
	ID          string     `json:"id"`           //唯一id
	ArtId		string 		`json:"art_id"`		 //	文章ID
	UserId      string     `json:"user_id"`      //用户id
	PartPersons []BaseUser `json:"part_persons"` //贡献者
	Version		int		`json:"version"`		//文章版本
	Token		string 	`json:"token"`			//文章令牌
	BaseArticle            //文章基本字段
	ArticleContent
	Comment
	DatetimeAware
}

type BaseArticle struct {
	Title   string `json:"title"`                          //标题
	Auth    string `json:"auth"`                           //主作者
	Sort    string `json:"sort"`                           //类别
	Content string `json:"content" gorm:"not null;unique"` //最新内容
}

type ArticleContent struct {
	SecondTitle string      `json:"second_title"` //副标题
	Photo       []BasePhoto `json:"photo"`        //图片
	Attachment  string      `json:"attachment"`   //附件
	Hot         string      `json:"hot"`          //热度
	Forward     string      `json:"forward"`      //转发数
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
	UserId string `json:"user_id"` //用户ID
	Token		string 	`json:"token"`			//文章令牌
	BaseArticle
	SecondTitle string      `json:"second_title"` //副标题
	Photo       []BasePhoto `json:"photo"`        //图片
	Attachment  string      `json:"attachment"`   //附件

}

type VerifyArticleRequest struct { //
	ID     string `json:"id"`      //文章ID
	UserId string `json:"user_id"` //用户id
	BaseArticle
	SecondTitle string      `json:"second_title"` //副标题
	Photo       []BasePhoto `json:"photo"`        //图片
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
	Content string `json:"content" gorm:"not null;unique"` //内容

}

func (Article) TableName() string {
	return TableArticle
}
