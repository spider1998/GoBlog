package entity

type YunXun struct {
	Sid        string `json:"sid"`        //用户的账号唯一标识
	Token      string `json:"token"`      //用户密钥Auth Token
	Appid      string `json:"appid"`      //应用分配标识
	Templateid string `json:"templateid"` //短信模板
	Param      string `json:"param"`      //参数
	Mobile     string `json:"mobile"`     //手机号
	Uid        string `json:"uid"`        //用户透传ID
}

type BindMobileRequest struct {
	ID     string `json:"id"`
	Mobile string `json:"mobile"`
	Mcode  string `json:"mcode"`
}
