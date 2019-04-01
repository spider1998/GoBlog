package routing

import (
	"Project/doit/handler/article"
	"Project/doit/handler/user"
	"Project/doit/routing/admin"
	"github.com/go-ozzo/ozzo-routing"
)

func ArticleRegisterRoutes(router *routing.RouteGroup) {
	var operatorHandler = new(admin.OperatorHandler)
	router.Get("/articles/<article_id>", article.GetArticle)         // 获取指定文章
	router.Get("/articles", article.GetArticles)                     // 获取全部文章
	router.Post("/add", article.AddArticle)                          // 创建文章
	router.Post("/verify", article.VerifyArticle)                    // 用户修改文章
	router.Post("/contribute", article.ContributeArticle)            //非用户修改文章
	router.Get("/<art_id>/comments", article.GetArticleComment)      //获取文章所有评论及回复
	router.Get("/sorts", operatorHandler.GetArticlesSorts)           //	获取文章分类
	router.Get("/articles/hot/top10", operatorHandler.GetArticleTop) // 获取文章排行前十

	router.Use(user.CheckSession)                                         // 检查用户登录状态信息
	router.Get("/arts/my", article.GetMyArticles)                         // 获取个人全部文章
	router.Get("/version/<article_id>", article.GetVersion)               // 获取文章所有版本
	router.Get("/view/<article_id>/<version>", article.GetVersionArticle) // 获取指定版本文章
	router.Post("/restore", article.RestoreVersionArticle)                // 恢复指定版本文章
	router.Delete("/<article_id>", article.DeleteArticle)                 // 删除文章
	router.Get("/likes/<likes_content>", article.QueryLikeArticles)       // 搜索相关文章
	router.Post("/update/<session_id>", article.UpdateArticle)            // 非用户用户修改文章

	router.Patch("/<article_id>/like", article.LikeOneArticle)    // 给文章点赞/取消点赞
	router.Get("/<article_id>/like", article.GetArticleLikeCount) // 获取文章点赞数

	router.Patch("/<article_id>", article.ForwardArticle)                //转发文章
	router.Patch("/authorization/forward", article.ForwardAuthorization) //文章转发授权
	router.Patch("/authorization/modify", article.ModifyAuthorization)   //文章修改授权

	router.Post("/<art_id>/comment", article.CommentArticle)      //评论文章
	router.Post("/comments/<com_id>/reply", article.CommentReply) //评论回复
}
