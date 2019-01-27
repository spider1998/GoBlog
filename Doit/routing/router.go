package routing

import (
	"Project/Doit/app"
	"Project/Doit/routing/article"
	"Project/Doit/routing/friend"
	"Project/Doit/routing/user"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/go-ozzo/ozzo-routing/slash"
	"net/http"
	"sync"
)

var (
	serverMutex sync.Mutex
	server      *http.Server
)

func Run() error {

	defer app.Logger.Info().Msg("http server terminated.")

	/*-----创建路由对象-----*/
	router := routing.New()
	router.Use(
		//添加跨域资源共享头
		cors.Handler(cors.Options{
			AllowOrigins:  "*",
			AllowHeaders:  "*",
			AllowMethods:  "*",
			ExposeHeaders: "X-Total-Count, X-Page, X-Page-Size",
		}),
		//记录请求日志输出		[192.168.183.1] [0.137ms] GET /version HTTP/1.1 200 61(字节)
		access.Logger(func(format string, a ...interface{}) {
			app.Logger.Info().Msgf(format, a...)
		}),
		//正确定向URL
		slash.Remover(http.StatusMovedPermanently),
		//设定MIME（扩展类型）格式
		content.TypeNegotiator(content.JSON),
		//错误处理
		errorHandler,
	)
	//版本信息路由校验
	router.Get("/version", func(c *routing.Context) error {
		var v struct {
			Version    string `json:"version"`
			CreateTime string `json:"create_time"`
		}
		v.Version = app.Version
		v.CreateTime = app.CreateTime
		return c.Write(v)
	})

	/*-----注册业务主路由-----*/
	app.Logger.Info().Msg("registering routes.")
	user.RegisterRoutes(router.Group("/user"))
	article.RegisterRoutes(router.Group("/article"))
	friend.RegisterRoutes(router.Group("/friend"))

	//遍历路由
	for _, route := range router.Routes() {
		app.Logger.Debug().Msgf("register route: \"%-6s -> %s\".", route.Method(), route.Path())
	}

	/*-----开始监听路由-----*/
	app.Logger.Info().Str("server_addr", app.Conf.HTTPAddr).Msg("run http server.")

	serverMutex.Lock()
	server = &http.Server{Addr: app.Conf.HTTPAddr, Handler: router}
	serverMutex.Unlock()
	err := http.ListenAndServe(app.Conf.HTTPAddr, router)
	if err != nil {
		if err == http.ErrServerClosed {
			return nil
		}
	}
	return err
}
