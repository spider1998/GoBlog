package article

import (
	"Project/Doit/handler/article"
	"Project/Doit/handler/user"
	"github.com/go-ozzo/ozzo-routing"
)

func RegisterRoutes(router *routing.RouteGroup) {
	router.Get("/<article_id>", article.GetArticle) // 获取指定文章
	router.Use(user.CheckSession)
	router.Get("/<article_id>/<version>",article.GetVersionArticle)				//获取指定版本文章
	router.Post("/add", article.AddArticle)       //创建文章
	router.Post("/verify", article.VerifyArticle) //用户修改文章
	router.Post("/update", article.UpdateArticle) //非用户用户修改文章
}
