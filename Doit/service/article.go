package service

import (
	"Project/Doit/app"
	"Project/Doit/code"
	"Project/Doit/entity"
	"Project/Doit/util"
	"github.com/go-ozzo/ozzo-dbx"
	v "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

var Article = ArticleService{}

type ArticleService struct{}

func (a *ArticleService) GetArticle(req string) (art entity.Article, err error) {
	err = app.DB.Select().Where(dbx.HashExp{"id": req}).One(&art)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

//创建文章
func (a *ArticleService) CreateArticle(req entity.CreateArticleRequest) (art entity.Article, err error) {

	err = v.ValidateStruct(&req,
		v.Field(&req.BaseArticle, v.Required),
	)
	if err != nil {
		return
	}
	art.BaseArticle = req.BaseArticle
	art.ID = uuid.New().String()
	art.UserId = req.UserId
	art.SecondTitle = req.SecondTitle
	art.Photo = req.Photo
	art.Attachment = req.Attachment
	art.CreateTime = util.DateTimeStd()
	art.UpdateTime = util.DateTimeStd()
	err = app.DB.Transactional(func(tx *dbx.Tx) error {
		err = tx.Model(&art).Insert()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if util.IsDBDuplicatedErr(err) {
			err = code.New(http.StatusConflict, code.CodeArticleExist)
			return
		}
		err = errors.Wrap(err, "fail to create article")
		return
	}
	return
}

func (a *ArticleService) VerifyArticle(req entity.VerifyArticleRequest) (art entity.Article, err error) {
	err = v.ValidateStruct(&req,
		v.Field(&req.BaseArticle, v.Required),
	)
	if err != nil {
		return
	}

	err = app.DB.Select().Where(dbx.HashExp{"id": req.ID}).One(&art)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	art.BaseArticle = req.BaseArticle
	art.Attachment = req.Attachment
	art.Photo = req.Photo
	art.SecondTitle = req.SecondTitle
	art.UpdateTime = util.DateTimeStd()

	err = app.DB.Model(&art).Update("Title", "Auth", "Sort", "Content", "Attachment", "Photo", "SecondTitle", "UpdateTime")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (a *ArticleService) UpdateArticle(req entity.UpdateArticleRequest, userId string) (status string, err error) {
	var art entity.Article
	err = v.ValidateStruct(&req,
		v.Field(&req.Content, v.Required),
	)
	if err != nil {
		return
	}
	err = app.DB.Select().Where(dbx.HashExp{"id": req.ID}).One(&art)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	if req.Content == art.Content {
		err = code.New(http.StatusBadRequest, code.CodeArticleNotChange)
		return
	}

}
