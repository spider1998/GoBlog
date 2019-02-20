package form

import "Project/doit/entity"

type CreateOperatorRequest struct {
	Name        string            `json:"name"`
	Password    string            `json:"password"`
	RealName    string            `json:"real_name"`
	Permissions map[string]string `json:"permissions"`
}

type UpdateOperatorRequest struct {
	ID          string            `json:"-"`
	RealName    string            `json:"real_name"`
	Permissions map[string]string `json:"permissions"`
}

type OperatorSignInRequest struct {
	Name         string `json:"name"`
	Password     string `json:"password"`
	CaptchaToken string `json:"captcha_token"`
	CaptchaCode  string `json:"captcha_code"`
}

type SetOperatorShortcutsRequest struct {
	ID        string        `json:"-"`
	Shortcuts []interface{} `json:"shortcuts"`
}

type QueryUserRequest struct {
	ID		string `json:"id"`
	State 	entity.UserState `json:"state"`
	Gender	entity.UserGender `json:"gender"`
	Oder 	string `json:"oder"`
}

type QueryArticleRequest struct {
	ID 		string `json:"id"`
	Sort 	string `json:"sort"`
}

type QueryArticleResponse struct {
	ID 		string `json:"id"`
	Sort 	string `json:"sort"`
	Auth 	string `json:"auth"`
	Title  	string `json:"title"`
	entity.DatetimeAware
}

type GetArticlesResponse struct {
	Count 	int `json:"count"`
	Arts 	[]QueryArticleResponse `json:"arts"`
}

type SiteStatisticResponse struct {
	UserCount		int `json:"user_count"`				//总用户数
	TodayRegister	int `json:"today_register"`			//今日注册用户数
	TodayArt		int `json:"today_art"`				//今日发布文章数
	ArtCount		int `json:"art_count"`				//总文章数
	ReadCount		int `json:"read_count"`				//总阅读量
}

type CreateArticleSortRequest struct {
	Name 	string `json:"name"`
	Sort 	string `json:"sort"`
}


