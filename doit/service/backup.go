package service

import (
	"Project/doit/app"
	"Project/doit/entity"
	"Project/doit/util"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

const (
	CronDBBackupKey = "db.backup"
)

var BackupApp = BackupAppService{
	cronReload: make(chan struct{}, 8),
}

type BackupAppService struct {
	cronReload chan struct{}
}

func (s *BackupAppService) ReloadCron() {
	s.cronReload <- struct{}{}
}

func (s *BackupAppService) RunCronLoop() {
	for {
		err := s.runCron()
		if err != nil {
			app.Logger.Error().Err(err).Msg("fail to run cron.")
			<-time.After(time.Second * 10)
		}
	}
}

func (s *BackupAppService) runCron() error {
	mFuncs := map[string]func() error{
		CronDBBackupKey: s.Backup,
	}

	app.Logger.Info().Msg("load cron.")

	var sysCrons []entity.SysCron
	err := app.DB.Select().All(&sysCrons)
	if err != nil {
		return errors.Wrap(err, "fail to find cron records.")
	}

	app.Logger.Info().Interface("crons", sysCrons).Msg("find cron from db.")

	c := cron.New()

	for i := range sysCrons {
		sysCron := sysCrons[i]
		err = c.AddFunc(sysCron.Spec, func() {
			app.Logger.Info().Str("key", sysCron.Key).Msg("add cron.")

			if _, ok := mFuncs[sysCron.Key]; !ok {
				app.Logger.Error().Err(err).Msg("fail to find cron func.")
				err = errors.New("fail to find cron func")
				return
			}

			err := mFuncs[sysCron.Key]()
			if err != nil {
				app.Logger.Error().Err(err).Msg("fail to run cron job.")
				return
			}

			sysCron.LastExecutedAt = util.DateTimeStd()
			err = app.DB.Model(&sysCron).Update("LastExecutedAt")
			if err != nil {
				app.Logger.Error().Err(err).Msg("fail to update last_executed_at of cron job.")
			}
		})
		if err != nil {
			return errors.Wrap(err, "fail to add cron func.")
		}

		app.Logger.Info().Str("key", sysCron.Key).Msg("added cron.")
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		c.Run()
	}()

	select {
	case <-s.cronReload:
		c.Stop()
	}

	wg.Wait()

	app.Logger.Info().Msg("stop cron.")

	return nil
}

