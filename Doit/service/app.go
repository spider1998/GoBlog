package service

import (
	"Project/Doit/app"
	"Project/Doit/code"
	"Project/Doit/entity"
	"Project/Doit/util"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"net/http"
	"strconv"
	"time"
)

var App = AppService{}

type AppService struct{}

func (a *AppService) CronRedis() {
	for {
		err := a.cronRun()
		if err != nil {
			app.Logger.Error().Err(err).Msg("fail to run cron !!!!!")
			timer1 := time.NewTicker(time.Minute * 1)
			for {
				select {
				case <-timer1.C:
					app.Logger.Info().Msg("try to reload ru cron !")
					a.CronRedis()
				}
			}
		}
	}
}

func (a *AppService) cronRun() (err error) {
	c := cron.New()
	spec := "@daily"
	c.AddFunc(spec, func() {
		vals, err := app.Redis.Cmd("KEYS", app.Conf.LikeRedis+"*").List()
		if err != nil {
			err = errors.Wrap(err, "fail to get  likes count from redis")
			return
		}
		for _, val := range vals {
			count, err := app.Redis.Cmd("SCARD", app.Conf.LikeRedis+":"+val).Int()
			if err != nil {
				return
			}
			var article entity.Article
			err = app.DB.Select().Where(dbx.HashExp{"art_id": val[7:]}).One(&article)
			if err != nil {
				if util.IsDBNotFound(err) {
					err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
					return
				}
				err = errors.WithStack(err)
				return
			}
			article.Hot = strconv.Itoa(count)
			err = app.DB.Model(&article).Update("Hot")
			if err != nil {
				err = errors.WithStack(err)
				return
			}
		}
	})

	c.Start()

	select {}
}
