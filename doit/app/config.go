package app

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	//基础服务配置
	Debug     bool   `json:"debug" default:"true"`      // 是否开启调试模式
	DebugHTTP bool   `json:"debug_http"`                //http调试
	HTTPAddr  string `json:"http_addr" default:":8081"` // HTTP 服务地址

	//数据库配置
	Mysql string `json:"mysql" default:"root:123456@tcp(192.168.35.193:3306)/doit?charset=utf8"` // mysql DSN
	Redis string `json:"redis" default:"192.168.35.193:6379"`
	LikeRedis string `json:"like_redis" default:"goblog"`

	//日志配置
	ConfPath       string `json:"conf_path" default:"."`        //日志文件路径
	PaginationPage int64  `json:"pagination_page" default:"1"`  //分页页数
	PaginationSize int64  `json:"pagination_size" default:"50"` //分页大小
	NSQD     string `json:"nsqd" default:"192.168.35.193:4150"`

	//邮件配置
	Email string `json:"email" default:"2387805574@qq.com"`		//服务器邮箱地址
	Epass string `json:"epass" default:"henuqnarpnucdjci"`		//邮箱密钥
	Etype string `json:"etype" default:"smtp.qq.com"`			//邮件服务器
	Eport int    `json:"eport" default:"587"`					//邮件服务端口

	//短信配置
	Msid    string `json:"msid" default:"116943babddda930dcd8802a7f6f5bd4"`        //用户唯一标识
	Mtoken  string `json:"mtoken" default:"28423c3bc2a1b63b4f432540e5b8cd96"`      //auth token
	Mappid  string `json:"mappid" default:"589910649b5347118abf1888f56a6071"`      //应用分配id
	Mcach   string `json:"mcach" default:"413802"`                                 //验证码模板id
	Mexpire string `json:"mexpire" default:"60"`                                   //验证码过期时间
	Maddr   string `json:"maddr" default:"https://open.ucpaas.com/ol/sms/sendsms"` //短信接口地址

	ContentSize	int `json:"content_size" default:"500"`								//版本控制分层数

	AttachmentPath	string `json:"attachment_path" default:"attachment/"`			//附件存储路径

}

//加载配置
func LoadConfig() (Config, error) {
	godotenv.Load()
	var config Config
	err := envconfig.Process("", &config)
	return config, err
}
