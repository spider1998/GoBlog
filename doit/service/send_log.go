package service

import (
	"encoding/json"
	"time"
	"github.com/nsqio/go-nsq"
	"Project/doit/form"
	"Project/doit/app"
)

var SLog = &SLogService{
	msgCH: make(chan form.CreateLogRequest, 1024),
}

type SLogService struct {
	msgCH chan form.CreateLogRequest
}

func (s *SLogService) SendLog(request form.CreateLogRequest) {
	app.Logger.Info().Msg("send log.")
	select {
	case s.msgCH <- request:
	default:
		app.Logger.Warn().Msg("log channel is full.")
	}
}

func (s *SLogService) Boot() error {
	go s.runLogSender()
	return nil
}

func (s *SLogService) runLogSender() {
	for {
		func() {
			app.Logger.Info().Msg("run log sender.")
			producer, err := nsq.NewProducer(app.Conf.NSQD, nsq.NewConfig())
			if err != nil {
				app.Logger.Error().Err(err).Msg("fail to create new dsq producer.")
				return
			}
			defer producer.Stop()
			producer.SetLogger(app.NewNSQLogger(app.Logger), nsq.LogLevelWarning)

			for msg := range s.msgCH {
				b, err := json.Marshal(msg)
				if err != nil {
					app.Logger.Error().Err(err).Msg("fail to marshal msg.")
					continue
				}
				err = producer.Publish("go-blog.log", b)
				if err != nil {
					app.Logger.Error().Err(err).Msg("fail to publish log message.")
					return
				}
			}
		}()
		time.Sleep(time.Second)
	}
}
