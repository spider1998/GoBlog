package article

import (
	"Project/Doit/entity"
	"Project/Doit/handler/session"
	"Project/Doit/service"
	"github.com/go-ozzo/ozzo-routing"
	"net/http"
	"crypto/sha1"
	"Project/Doit/code"
	"Project/Doit/app"
	"encoding/hex"
	"strconv"
)

//获取最新版文章
func GetArticle(c *routing.Context) error {
	req := c.Param("article_id")
	article, err := service.Article.GetArticle(req)
	if err != nil {
		return err
	}
	return c.Write(article)
}

//获取指定版本文章
func GetVersionArticle(c *routing.Context) error {
	ver := c.Param("version")
	artId := c.Param("article_id")
	version,err := strconv.Atoi(ver)
	if err != nil{
		return err
	}
	article,err := service.Article.GetVersionArticle(version,artId)
	if err != nil {
		return err
	}
	return c.Write(article)
}

//恢复历史版本
func RestoreVersionArticle(c *routing.Context) error {
	var req entity.RestoreArticleRequest
	err := c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	req.UserId = session.GetUserSession(c).ID

	article,err := service.Article.RestoreVersionArticle(req)
	if err != nil {
		return err
	}
	return c.Write(article)

}

//添加文章
func AddArticle(c *routing.Context) error {
	var request entity.CreateArticleRequest
	err := c.Read(&request)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	request.UserId = session.GetUserSession(c).ID
	art, err := service.Article.CreateArticle(request)
	if err != nil {
		return err
	}
	if err = c.Write(art);err !=nil{
		return err
	}
	var hashContent entity.Content
	var hc string = ""
	var de string = ""
	for i := 0;;i+=app.Conf.ContentSize{
		if i >= len(art.Content){
			break
		}
		if i + app.Conf.ContentSize > len(art.Content){
			de = art.Content[i:]

		}else {
			de = art.Content[i:i+app.Conf.ContentSize]
		}
		hashContent.VersionHash = art.Version								//片段版本
		hashContent.Detail = de											//详细内容
		hashContent.HeadUuid = hc										//头标识
		hs := sha1.Sum([]byte(hashContent.Detail))
		hc = hex.EncodeToString(hs[:])
		hashContent.TailUuid = hc										//尾标识
		hashContent.Changed = false										//改动标识
		hashContent.UserId = art.UserId									//用户ID
		hashContent.ArtId = art.ArtId									//文章ID
		err := service.Article.SaveArtBlock(hashContent)
		if err != nil{
			return err
		}
	}
	return nil
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
	return c.Write(response)
}
