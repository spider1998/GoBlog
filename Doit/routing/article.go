package routing

import (
	"Project/Doit/handler/article"
	"github.com/go-ozzo/ozzo-routing"
	"Project/Doit/handler/user"
)

func ArticleRegisterRoutes(router *routing.RouteGroup) {
	router.Get("/articles/<article_id>", article.GetArticle)								// 获取指定文章
	router.Get("/articles", article.GetArticles)								// 获取指定文章
	router.Post("/add", article.AddArticle)       								// 创建文章
	router.Post("/verify", article.VerifyArticle) 								// 用户修改文章


	router.Use(user.CheckSession)													// 检查用户登录状态信息
	router.Get("/version/<article_id>",article.GetVersion)						// 获取文章所有版本
	router.Get("/view/<article_id>/<version>",article.GetVersionArticle)		// 获取指定版本文章
//	router.Post("/restore",article.RestoreVersionArticle)							// 恢复指定版本文章
	router.Delete("/<article_id>",article.DeleteArticle)						// 删除文章
	router.Get("/likes/<likes_content>",article.QueryLikeArticles)				// 搜索相关文章
	router.Post("/update", article.UpdateArticle) 								// 非用户用户修改文章


	router.Patch("/<article_id>/like",article.LikeOneArticle)					// 给文章点赞/取消点赞
	router.Get("/<article_id>/like",article.GetArticleLikeCount)				// 获取文章点赞数

	router.Patch("/<article_id>",article.ForwardArticle)						//转发文章
}
