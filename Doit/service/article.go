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
	"github.com/mediocregopher/radix.v2/redis"
	"fmt"
	"github.com/rs/xid"
)


var Article = ArticleService{
	sessionExp: 86400,
}

type ArticleService struct{
	sessionExp int
}

//获取最新版本文章
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
	art.Read += 1
	err = app.DB.Model(&art).Update("Read")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//获取最新版本文章
func (a *ArticleService) GetArticles() (arts []entity.Article, err error) {
	err = app.DB.Select().OrderBy("create_time desc").All(&arts)
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


//获取文章所有版本，返回版本列表
func (a *ArticleService) GetVersion(req string) (version []int,err error) {
	var con []entity.ArticleVersion
	err = app.DB.Select().Where(dbx.HashExp{"art_id": req}).OrderBy("version desc").All(&con)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	for _,c := range con{
		version = append(version,c.Version )
	}
	return
}

//获取历史版本文章
func (a *ArticleService) GetVersionArticle(version int,artId string) (art entity.Article,err error) {
	var con entity.ArticleVersion
	err = app.DB.Select().Where(dbx.HashExp{"art_id": artId}).
		AndWhere(dbx.NewExp("version={:ver}", dbx.Params{"ver": version})).One(&con)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	err = app.DB.Select().Where(dbx.HashExp{"id": artId}).One(&art)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	art.Title = con.Title
	art.SecondTitle = con.SecondTitle
	art.Content = con.Content
	art.Version = version
	art.Photo = con.Photo
	return
}

//删除高版本文章缓存
func (a *ArticleService)DeleteMaxArticle(version int) (err error) {
	var cons []entity.ArticleVersion
	err = app.DB.Delete("article_version",dbx.NewExp("version>{:ver}", dbx.Params{"ver": version})).All(&cons)
	if err != nil{
		return err
	}
	return
}

//查询相关标题文章
func (a *ArticleService) QueryLikeArticles(content string) (arts []entity.Article,err error) {
	err = app.DB.Select().Where(dbx.NewExp("title like %{:con}%", dbx.Params{"con": content})).Where(dbx.NewExp("second_title like %{:con}%", dbx.Params{"con": content})).All(&arts)
	if err = DbErrorHandler(err, false); err != nil {
		return
	}
	return
}

//删除文章
func (a *ArticleService)DeleteArticle(articleID,userID string) (err error) {
	var art entity.Article
	err = app.DB.Select().Where(dbx.HashExp{"id": articleID}).One(&art)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	if userID != art.UserId{
		fmt.Println("==========")
		err = code.New(http.StatusBadRequest, code.CodeDenied)
		return
	}
	err = app.DB.Model(&art).Delete()
	if err != nil{
		return err
	}
	return
}

//恢复历史版本
func (a *ArticleService)RestoreVersionArticle(req entity.RestoreArticleRequest) (art entity.Article,err error)  {
	err = v.ValidateStruct(&req,
		v.Field(&req.Version, v.Required),
		v.Field(&req.ArtId, v.Required),
	)
	if err != nil {
		return
	}
	//查询指定文章
	err = app.DB.Select().Where(dbx.HashExp{"id": req.ArtId}).One(&art)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	if req.UserId != art.UserId{
		err = code.New(http.StatusBadRequest, code.CodeDenied)
		return
	}
	//查询指定版本文章
	var verArt entity.ArticleVersion
	err = app.DB.Select().Where(dbx.HashExp{"art_id": req.ArtId}).One(&verArt)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	art.Title = verArt.Title
	art.SecondTitle = verArt.SecondTitle
	art.ModifyType = verArt.ModifyType
	art.Sort = verArt.Sort
	art.Content = verArt.Content
	art.Photo = verArt.Photo
	art.Attachment = verArt.Attachment
	art.Version = req.Version
	art.UpdateTime = util.DateTimeStd()

	err = app.DB.Model(&art).Update("Content", "Version", "UpdateTime",
		"Title","SecondTitle","ModifyType","Sort","Photo","Attachment")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
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
	for len(con)!=0{
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
	u,err := User.CheckSession(req.UserId)
	if err != nil{
		return
	}
	var user entity.User
	err = app.DB.Select().Where(dbx.HashExp{"id": u.ID}).One(&user)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusNotFound, code.CodeUserAccessSessionInvalid).Err("account not found.")
		}
		err = errors.Wrap(err, "fail to find user")
		return
	}
	art.Content = req.Content
	art.Title = req.Title
	art.Auth = user.Name
	art.Sort = req.Sort
	art.Version = 1
	art.ID = uuid.New().String()
	art.UserId = u.ID
	art.SecondTitle = req.SecondTitle
	art.Photo = req.Photo
	art.Attachment = req.Attachment
	art.Hot = 0
	art.PartPersons += ","+u.ID
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


func (a *ArticleService) VerifyArticle(req entity.VerifyArticleRequest) (err error) {
		var art entity.Article
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
	var artV entity.ArticleVersion
	artV.ArtID = art.ID
	artV.UserId = art.UserId
	artV.Version = art.Version
	artV.ModifyType= art.ModifyType
	artV.BaseArticle= art.BaseArticle
	artV.ArticleContent= art.ArticleContent
	artV.Comment= art.Comment
	artV.UpdateTime= util.DateTimeStd()
	artV.ArtID = art.ID
	artV.ID = xid.New().String()
	err = app.DB.Transactional(func(tx *dbx.Tx) error {
		err = tx.Model(&artV).Insert()
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
	//更新文章
	art.Version += 1
	art.BaseArticle = req.BaseArticle
	art.Attachment = req.Attachment
	art.Photo = req.Photo
	art.SecondTitle = req.SecondTitle
	art.UpdateTime = util.DateTimeStd()
	err = app.DB.Model(&art).Update("Title", "Auth", "Sort", "Content", "Attachment", "Photo", "SecondTitle", "UpdateTime","Version")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//非用户修改文章
func (a *ArticleService) UpdateArticle(req entity.UpdateArticleRequest, userId string) (art entity.Article, err error) {
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
	return
}

//文章转发授权
func (a *ArticleService)ForwardAuthorazation(req entity.ArticleAuthorazation)(err error) {
	var art entity.ArticleForward
	err = app.DB.Select().Where(dbx.HashExp{"id": req.RecordID}).One(&art)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	if req.State == 3{
		art.Status = entity.StateForwardFinished
	}
	if req.State == 2{
		art.Status = entity.StateForwardRefused
	}
	err = app.DB.Model(&art).Update("Satus")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//文章转发
func (a *ArticleService) ForwardArticle(req entity.ArticleForwardRequest)(err error) {
	err = v.ValidateStruct(&req,
		v.Field(&req.Reason, v.Required),
	)
	if err != nil {
		return
	}
	var art entity.Article
	err = app.DB.Select().Where(dbx.HashExp{"id": req.ArtID}).One(&art)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	var res entity.ArticleForward
	res.ID = xid.New().String()
	res.Reason = req.Reason
	res.ArtID = req.ArtID
	res.ForwardID = req.UserID
	res.AuthID = art.UserId
	res.Status = entity.StateForwardWaite
	res.CreateTime = util.DateTimeStd()
	res.UpdateTime = util.DateTimeStd()
	err = app.DB.Transactional(func(tx *dbx.Tx) error {
		err = tx.Model(&res).Insert()
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
		err = errors.Wrap(err, "fail to create forward article")
		return
	}
	return

}

//获取文章点赞次数
func (a *ArticleService) GetArticleLikeCount(artID string) (count int,err error) {
	val,err := app.Redis.Cmd("EXISTS", app.Conf.LikeRedis+":"+artID).Int()
	if err != nil {
		if err == redis.ErrRespNil {
			err = code.New(http.StatusForbidden, code.CodeUserAccessSessionInvalid).Err("record session not found.")
			return
		}
		err = errors.Wrap(err, "fail to get  likes count from redis")
		return
	}
	if val == 1{
		count,err = app.Redis.Cmd("SCARD",app.Conf.LikeRedis+":"+artID).Int()
		if err != nil{
			return
		}
	}else {
		var article entity.Article
		err1 := app.DB.Select("hot").Where(dbx.HashExp{"id": artID}).One(&article)
		if err1 != nil {
			if util.IsDBNotFound(err) {
				err1 = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
				err = err1
				return
			}
			err = errors.WithStack(err)
			return
		}
		count = article.Hot
		if err != nil{
			return
		}
	}
	return
}

//文章点赞/取消带点赞
func (a *ArticleService) LikeOneArticle(articleID,userID string) (err error) {
	val,err := app.Redis.Cmd("SISMEMBER", app.Conf.LikeRedis+":"+articleID,userID).Int()
	if err != nil {
		if err == redis.ErrRespNil {
			err = code.New(http.StatusForbidden, code.CodeUserAccessSessionInvalid).Err("record session not found.")
			return
		}
		err = errors.Wrap(err, "fail to get email code from redis")
		return
	}
	if val == 1{
		err1 := app.Redis.Cmd("SREM", app.Conf.LikeRedis+":"+articleID,userID).Err
		if err1 != nil {
			if err1 == redis.ErrRespNil {
				err1 = code.New(http.StatusForbidden, code.CodeUserAccessSessionInvalid).Err("record session not found.")
				err = err1
				return
			}
			err = errors.Wrap(err, "fail to delete like members from redis")
			return
		}
	}else {
		err = app.Redis.Cmd("SADD", app.Conf.LikeRedis+":"+articleID, userID).Err
		if err != nil {
			err = errors.Wrap(err, "fail to set like members redis")
			return
		}
	}
	return
}
