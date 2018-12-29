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
	"crypto/sha1"
	"github.com/gobuffalo/packr/v2/file/resolver/encoding/hex"
)

var Article = ArticleService{}

type ArticleService struct{}

//获取最新版本文章
func (a *ArticleService) GetArticle(req string) (art entity.Article, err error) {
	err = app.DB.Select().Where(dbx.HashExp{"art_id": req}).One(&art)
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

//获取历史版本文章
func (a *ArticleService) GetVersionArticle(version int,artId string) (art entity.Article,err error) {
	var con []entity.Content
	err = app.DB.Select().Where(dbx.HashExp{"art_id": artId}).
		AndWhere(dbx.NewExp("version<={:ver}", dbx.Params{"ver": version})).
			AndWhere(dbx.HashExp{"changed": false}).All(&con)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}

	err = app.DB.Select().Where(dbx.HashExp{"art_id": artId}).One(&art)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	token := art.Token
	content := LinkBlock(con,token)
	art.Content = content
	art.Version = version
	return
}

//链接文章块
func LinkBlock(con []entity.Content,token string) (string)  {
	var content string
	hs := sha1.Sum([]byte(token))
	node := hex.EncodeToString(hs[:])
	co := con[0]
	con = con[1:]
	content = node + co.Detail
	for len(con)==0{
		for j,c := range con{
			if c.HeadUuid == co.TailUuid{
				co = c
				content = content + node + co.Detail
				con = append(con[:j],con[j+1:]...)
				break
			}
			if c.TailUuid == co.HeadUuid{
				co = c
				content = node + co.Detail+ content
				con = append(con[:j],con[j+1:]...)
				break
			}
		}
	}
	return content

}


//创建文章
func (a *ArticleService) CreateArticle(req entity.CreateArticleRequest) (art entity.Article, err error) {

	err = v.ValidateStruct(&req,
		v.Field(&req.BaseArticle, v.Required),
	)
	if err != nil {
		return
	}
	art.Token = req.Token
	art.Title = req.Title
	art.Auth = req.Auth
	art.Sort = req.Sort
	art.Version = 1
	art.ArtId = uuid.New().String()
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

//存储文章区块
func (a *ArticleService)SaveArtBlock(req entity.Content) (err error) {
	err = v.ValidateStruct(&req,
		v.Field(&req.Detail, v.Required),
	)
	if err != nil {
		return
	}
	//用户自己修改，无需审核
	err = app.DB.Transactional(func(tx *dbx.Tx) error {
		err = tx.Model(&req).Insert()
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
		err = errors.Wrap(err, "fail to create article block")
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
