package routing

import (
	"Project/Doit/app"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/go-ozzo/ozzo-routing/slash"
	"net/http"
	"sync"
	"os"
	"regexp"
	"io"
	"github.com/go-ozzo/ozzo-routing/file"
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
			ExposeHeaders: "X-Page-Total, X-Page, X-Page-Size",
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
	router.Get("/static/*", file.Server(file.PathMap{
		"/static/": "../",
	}))
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
	UserRegisterRoutes(router.Group("/user"))
	ArticleRegisterRoutes(router.Group("/article"))
	FriendRegisterRoutes(router.Group("/friend"))


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

func fileStatic(w http.ResponseWriter, r *http.Request) {
	if ok, _ := regexp.MatchString("/static/", r.URL.String()); ok {
		StaticServer(w, r)
		return
	}
	io.WriteString(w, "hello world")
}

func StaticServer(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		app.Logger.Fatal().Err(err)
	}
	http.StripPrefix("/static/",
		http.FileServer(http.Dir(wd))).ServeHTTP(w, r)
}