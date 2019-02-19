package service

import (
	"encoding/json"
	"time"

	v "github.com/go-ozzo/ozzo-validation"
	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"Project/doit/entity"
	"Project/doit/app"
	"github.com/go-ozzo/ozzo-dbx"
	"Project/doit/util"
	"Project/doit/code"
	"net/http"
	"Project/doit/form"
)

var Log = &LogService{}


type LogService struct{}

func (s *LogService) Boot() error {
	go new(LogCollector).RunLoop()
	return nil
}

func (s *LogService) LogOperator(operator entity.Operator, system, action, remark, ip string, ext ...util.M) (log entity.Log, err error) {
	return s.Log(entity.LogUserTypeOperator, operator.ID, operator.Name, system, action, remark, ip, ext...)
}

func (s *LogService) Log(userType entity.LogUserType, userID, userName, system, action, remark, ip string, ext ...util.M) (log entity.Log, err error) {
	log.ID = xid.New().String()
	log.UserType = userType
	log.UserID = userID
	log.UserName = userName
	log.System = system
	log.Action = action
	log.Remark = remark
	log.IP = ip
	if len(ext) > 0 {
		log.Ext = ext[0]
	}

	err = v.ValidateStruct(&log,
		v.Field(&log.UserType, v.Required, v.In(entity.LogUserTypeUser, entity.LogUserTypeOperator)),
		v.Field(&log.UserID, v.Required),
		v.Field(&log.UserName, v.Required),
		v.Field(&log.System, v.Required, v.RuneLength(1, 32)),
		v.Field(&log.Action, v.Required, v.RuneLength(1, 32)),
		v.Field(&log.Remark, v.Required),
		v.Field(&log.IP, v.Required),
	)
	if err != nil {
		app.Logger.Error().Interface("log", log).Err(err).Msg("fail to validate log.")
		err = errors.WithStack(err)
		return
	}

	err = app.DB.Transactional(func(tx *dbx.Tx) error {
		err = tx.Model(&log).Insert()
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


func (s *LogService) CountLogs(cond form.QueryLogsCond) (n int, err error) {
	sess := app.DB.Select("count(*)").From(entity.TableLog)
	if cond.UserType > 0 {
		sess.AndWhere(dbx.HashExp{"user_type": cond.UserType})
	}
	if cond.Remark != "" {
		sess.AndWhere(dbx.Like("remark",cond.Remark))
	}
	if cond.FromTime != "" {
		sess.AndWhere(dbx.NewExp("create_time>={:ct}",dbx.Params{"ct":cond.FromTime}))
	}
	if cond.ToTime != "" {
		sess.AndWhere(dbx.NewExp("create_time<={:ct}",dbx.Params{"ct":cond.ToTime}))
	}

	//记录总数
	var cnt int
	err = sess.Row(&cnt)
	if err != nil {
		err = errors.Wrap(err, "fail to query devices.")
		return
	}
	return int(cnt), nil
}

func (s *LogService) QueryLogs(offset, limit int, cond form.QueryLogsCond) (logs []entity.Log, err error) {
	sess := app.DB.Select("count(*)").From(entity.TableLog)
	if cond.UserType > 0 {
		sess.AndWhere(dbx.HashExp{"user_type": cond.UserType})
	}
	if cond.Remark != "" {
		sess.AndWhere(dbx.Like("remark",cond.Remark))
	}
	if cond.FromTime != "" {
		sess.AndWhere(dbx.NewExp("create_time>={:ct}",dbx.Params{"ct":cond.FromTime}))
	}
	if cond.ToTime != "" {
		sess.AndWhere(dbx.NewExp("create_time<={:ct}",dbx.Params{"ct":cond.ToTime}))
	}
	err = sess.OrderBy("create_time desc").Limit(int64(limit)).Offset(int64(offset)).All(&logs)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if logs == nil {
		logs = make([]entity.Log, 0)
	}
	return
}


type LogCollector struct{}

func (c *LogCollector) RunLoop() {
	for {
		c.collect()
		<-time.After(time.Second)
	}
}

func (c *LogCollector) collect() {
	app.Logger.Info().Msg("start nsq runner.")
	consumer, err := nsq.NewConsumer("go-blog.log", app.System, nsq.NewConfig())
	if err != nil {
		app.Logger.Error().Err(err).Msg("fail to create nsq consumer.")
		return
	}
	consumer.SetLogger(app.NewNSQLogger(app.Logger), nsq.LogLevelWarning)
	consumer.AddHandler(c)
	err = consumer.ConnectToNSQD(app.Conf.NSQD)
	if err != nil {
		app.Logger.Error().Err(err).Msg("fail to connect to nsqd.")
		return
	}
	<-consumer.StopChan
	app.Logger.Info().Msg("nsq runner stopped.")
}

func (c *LogCollector) HandleMessage(msg *nsq.Message) error {
	app.Logger.Debug().Str("body", string(msg.Body)).Msg("received nsq log message.")
	var request form.CreateLogRequest
	err := json.Unmarshal(msg.Body, &request)
	if err != nil {
		app.Logger.Error().Err(err).Msg("fail to decode log request.")
		return nil
	}
	switch request.UserType {
	case entity.LogUserTypeOperator:
		op, err := Operator.CheckToken(request.Token)
		if err != nil {
			app.Logger.Error().Err(err).Msg("fail to check operator token.")
			return nil
		}
		log, err := Log.LogOperator(op, request.System, request.Action, request.Remark, request.IP, request.Ext)
		if err != nil {
			app.Logger.Error().Err(err).Msg("fail to add operator log。")
			return nil
		}
		app.Logger.Info().Str("log_id", log.ID).Msg("add operator log successfully.")
		return nil
	}

	app.Logger.Info().Int("user_type", int(request.UserType)).Msg("unsupported log type.")
	return nil
}
