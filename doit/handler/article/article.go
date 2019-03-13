package article

import (
	"Project/doit/app"
	"Project/doit/code"
	"Project/doit/entity"
	"Project/doit/form"
	"Project/doit/handler/session"
	"Project/doit/service"
	"Project/doit/util"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

//获取文章
func GetArticle(c *routing.Context) error {
	req := c.Param("article_id")
	article, err := service.Article.GetArticle(req)
	if err != nil {
		return err
	}
	return c.Write(article)
}

//获取整站全部文章
func GetArticles(c *routing.Context) error {
	article, err := service.Article.GetArticles()
	if err != nil {
		return err
	}
	return c.Write(article)
}

//获取个人文章
func GetMyArticles(c *routing.Context) error {
	userID := session.GetUserSession(c).ID
	article, err := service.Article.GetMyArticles(userID)
	if err != nil {
		return err
	}
	return c.Write(article)
}

//获取历史版本
func GetVersion(c *routing.Context) error {
	req := c.Param("article_id")
	version, err := service.Article.GetVersion(req)
	if err != nil {
		return err
	}
	return c.Write(version)
}

//获取指定版本文章
func GetVersionArticle(c *routing.Context) error {
	ver := c.Param("version")
	artId := c.Param("article_id")
	version, err := strconv.Atoi(ver)
	if err != nil {
		return err
	}
	article, err := service.Article.GetVersionArticle(version, artId)
	if err != nil {
		return err
	}
	return c.Write(article)
}

//删除文章
func DeleteArticle(c *routing.Context) error {
	articleID := c.Param("article_id")
	userID := session.GetUserSession(c).ID
	err := service.Article.DeleteArticle(articleID, userID)
	if err != nil {
		return err
	}
	return c.Write(http.StatusOK)
}

//查询相关标题文章
func QueryLikeArticles(c *routing.Context) error {
	content := c.Param("likes_content")
	response, err := service.Article.QueryLikeArticles(content)
	if err != nil {
		return err
	}

	if len(response) == 0 {
		response = []entity.Article{}
	}
	pager := util.GetPaginatedListFromRequest(c, len(response))
	if pager.Offset()+pager.Limit() <= pager.TotalCount {
		return c.Write(response[pager.Offset() : pager.Offset()+pager.Limit()])
	} else {
		return c.Write(response[pager.Offset():pager.TotalCount])
	}
}

//恢复历史版本(同时删除大于该版本的所有版本)
func RestoreVersionArticle(c *routing.Context) error {
	var req entity.RestoreArticleRequest
	err := c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	req.UserId = session.GetUserSession(c).ID

	article, err := service.Article.RestoreVersionArticle(req)
	if err != nil {
		return err
	}
	err = service.Article.DeleteMaxArticle(req.Version)
	if err != nil {
		return err
	}
	return c.Write(article)

}

//保存文章变动区块
func SaveVerified(art entity.Article) (err error) {
	var hashContent entity.Content
	var hc string = ""
	var de string = ""
	//拆分文章重新结合为带标识的文章块
	for i := 0; ; i += app.Conf.ContentSize {
		if i >= len(art.Content) {
			break
		}
		if i+app.Conf.ContentSize > len(art.Content) {
			de = art.Content[i:]

		} else {
			de = art.Content[i : i+app.Conf.ContentSize]
		}
		hashContent.Version = art.Version //片段版本
		hashContent.Detail = de           //详细内容
		hashContent.HeadUuid = hc         //头标识
		hs := sha1.Sum([]byte(hashContent.Detail))
		hc = hex.EncodeToString(hs[:])
		hashContent.TailUuid = hc       //尾标识
		hashContent.Changed = false     //改动标识
		hashContent.UserId = art.UserId //用户ID
		hashContent.ArtId = art.ID      //文章ID
		err := service.Article.SaveArtBlock(hashContent)
		if err != nil {
			return err
		}
	}
	return nil
}

//添加文章
func AddArticle(c *routing.Context) error {
	var req entity.CreateArticleRequest
	req.UserId = c.Form("user")
	req.Title = c.Form("title")
	req.SecondTitle = c.Form("second_title")
	modify := c.Form("modify_type")
	if modify == "1" {
		req.ModifyType = entity.ModifyTypeAble
	} else {
		req.ModifyType = entity.ModifyTypeEnable
	}
	req.Sort = c.Form("sort")
	req.Content = c.Form("content")

	imgF, imgH, err := c.Request.FormFile("bacc")
	if err != nil {
		app.Logger.Info().Msg("no img")
	}
	//保存背景图
	if imgF == nil {
		fmt.Println("no img")
	} else {
		imgPath, err := saveFile(imgF, imgH)
		if err != nil {
			return err
		}
		defer imgF.Close()
		req.Photo = imgPath
	}
	//保存附件
	testF, testH, err := c.Request.FormFile("attach")
	if err != nil {
		app.Logger.Info().Msg("no attachment")
	}
	if testF == nil {
		fmt.Println("no img")
	} else {
		testPath, err := saveFile(testF, testH)
		if err != nil {
			return err
		}
		defer testF.Close()
		req.Attachment = testPath
	}
	art, err := service.Article.CreateArticle(req)
	if err != nil {
		return err
	}
	return c.Write(art.ID)
}

func saveFile(file multipart.File, head *multipart.FileHeader) (path string, err error) {
	path = service.User.SaveAttachment(head)
	if _, err1 := os.Stat(path); err1 != nil {
		err1 := os.MkdirAll(path, 0711)
		if err1 != nil {
			err = err1
			return
		}
	}

	fW, err := os.Create(path + head.Filename)
	if err != nil {
		return
	}
	defer fW.Close()
	io.Copy(fW, file)
	path = path + head.Filename
	return
}

//用户修改文章
func VerifyArticle(c *routing.Context) error {

	var req entity.VerifyArticleRequest
	req.ID = c.Form("artID")
	req.UserId = c.Form("user")

	u, err := service.User.CheckSession(req.UserId)
	if err != nil {
		return err
	}
	var art entity.Article
	err = app.DB.Select().Where(dbx.HashExp{"id": req.ID}).One(&art)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return err
		}
		err = errors.WithStack(err)
		return err
	}
	if u.ID != art.UserId {
		err = code.New(http.StatusBadRequest, code.CodeIDNotAllowed)
		return err
	}

	req.Title = c.Form("title")
	req.SecondTitle = c.Form("second_title")
	modify := c.Form("modify_type")
	if modify == "1" {
		req.ModifyType = entity.ModifyTypeAble
	} else {
		req.ModifyType = entity.ModifyTypeEnable
	}
	req.Sort = c.Form("sort")
	req.Content = c.Form("content")
	imgF, imgH, err := c.Request.FormFile("bacc")
	if err != nil {
		app.Logger.Info().Msg("no img")
	}
	//保存背景图
	if imgF == nil {
		fmt.Println("no img")
	} else {
		imgPath, err := saveFile(imgF, imgH)
		if err != nil {
			return err
		}
		defer imgF.Close()
		req.Photo = imgPath
	}
	//保存附件
	testF, testH, err := c.Request.FormFile("attach")
	if err != nil {
		app.Logger.Info().Msg("no attachment")
	}
	if testF == nil {
		fmt.Println("no img")
	} else {
		testPath, err := saveFile(testF, testH)
		if err != nil {
			return err
		}
		defer testF.Close()
		req.Attachment = testPath
	}
	err = service.Article.VerifyArticle(req)
	if err != nil {
		return err
	}
	return c.Write(http.StatusOK)

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

// 点赞/取消点赞操作
func LikeOneArticle(c *routing.Context) error {
	articleID := c.Param("article_id")
	userID := session.GetUserSession(c).ID
	err := service.Article.LikeOneArticle(articleID, userID)
	if err != nil {
		return err
	}
	return c.Write(http.StatusOK)

}

//获取文章点赞数量
func GetArticleLikeCount(c *routing.Context) error {
	artID := c.Param("article_id")
	count, err := service.Article.GetArticleLikeCount(artID)
	if err != nil {
		return err
	}
	return c.Write(count)
}

//文章转发
func ForwardArticle(c *routing.Context) error {
	var req entity.ArticleForwardRequest
	err := c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	authID, err := service.Article.ForwardArticle(req)
	if err != nil {
		return err
	}
	//发送转发授权消息
	err = service.Message.Create(authID, app.Conf.BlogForward, req.ArtID, req.Reason, req.UserID)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeCreateMessageError)
	}
	return c.Write(http.StatusOK)
}

//文章转发授权
func ForwardAuthorization(c *routing.Context) error {
	var req entity.ArticleAuthorization
	err := c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	content, authID, err := service.Article.ForwardAuthorization(req)
	if err != nil {
		return err
	}
	//发送授权消息
	err = service.Message.Create(req.UserID, app.Conf.BlogForward, req.ArtID, content, authID)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeCreateMessageError)
	}
	return c.Write(http.StatusOK)
}

//文章评论
func CommentArticle(c *routing.Context) error {
	var req form.CommentArticleRequest
	err := c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	req.ArtID = c.Param("art_id")
	req.UserID = session.GetUserSession(c).ID
	req.Name = session.GetUserSession(c).Name
	err = service.Article.CommentArticle(req)
	if err != nil {
		return err
	}
	return c.Write(http.StatusOK)
}

//评论回复
func CommentReply(c *routing.Context) error {
	var req form.CommentReplyRequest
	err := c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}
	req.ComID = c.Param("com_id")
	req.UserID = session.GetUserSession(c).ID
	req.Name = session.GetUserSession(c).Name
	err = service.Article.CommentReply(req)
	if err != nil {
		return err
	}
	return c.Write(http.StatusOK)
}

//获取所有评论及回复
func GetArticleComment(c *routing.Context) error {
	artID := c.Param("art_id")
	if artID == "" {
		return code.New(http.StatusBadRequest, code.CodeBadRequest)
	}
	response, err := service.Article.GetArticleComment(artID)
	if err != nil {
		return err
	}
	if len(response) == 0 {
		response = []form.ArticleCommentResponse{}
	}
	pager := util.GetPaginatedListFromRequest(c, len(response))
	if pager.Offset()+pager.Limit() <= pager.TotalCount {
		return c.Write(response[pager.Offset() : pager.Offset()+pager.Limit()])
	} else {
		return c.Write(response[pager.Offset():pager.TotalCount])
	}

}
