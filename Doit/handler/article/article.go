package article

import (
	"DBInsert/code"
	"Project/Doit/entity"
	"Project/Doit/handler/session"
	"Project/Doit/service"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/mediocregopher/radix.v2/redis"
	"net/http"
)

func GetArticle(c *routing.Context) error {
	req := c.Param("article_id")
	article, err := service.Article.GetArticle(req)
	if err != nil {
		return err
	}
	return c.Write(article)
}

func AddArticle(c *routing.Context) error {
	var request entity.CreateArticleRequest
	err := c.Read(&request)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	request.UserId = session.GetUserSession(c).ID
	respons, err := service.Article.CreateArticle(request)
	if err != nil {
		return err
	}
	return c.Write(respons)

}

//用户修改文章
func VerifyArticle(c *routing.Context) error {
	var verify entity.VerifyArticleRequest
	id := session.GetUserSession(c).ID
	err := c.Read(&verify)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	if verify.UserId != id {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	response, err := service.Article.VerifyArticle(verify)
	if err != nil {
		return err
	}
	return c.Write(response)

}

//非用户修改文章
func UpdateArticle(c *routing.Context) error {
	var request entity.UpdateArticleRequest
	err := c.Read(&request)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	userId := session.GetUserSession(c).ID
	response, err := service.Article.UpdateArticle(request, userId)
	if err != nil {
		return err
	}

}
