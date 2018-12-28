package app

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/rs/zerolog"
	"time"
)

type RedisBox struct {
	*pool.Pool
	logger zerolog.Logger
}

//创建自定义redis
func LoadRedis(addr string, logger zerolog.Logger) (box RedisBox, err error) {
	//指定一个DialFunc，在池新连接时使用
	p, err := pool.NewCustom("tcp", addr, 10, func(network, addr string) (*redis.Client, error) {
		//限时连接
		return redis.DialTimeout(network, addr, time.Second*3)
	})
	if err != nil {
		return
	}
	box.Pool = p
	box.logger = logger
	return
}

//在池中获取一个客户端执行命令操作，完成后放回池中
func (r RedisBox) Cmd(cmd string, args ...interface{}) *redis.Resp {
	t := time.Now()
	resp := r.Pool.Cmd(cmd, args...)
	//日志记录详细命令
	r.logger.Debug().Str("cmd", cmd).Dur("elapsed", time.Now().Sub(t)).Interface("args", args).Str("resp", resp.String()).Msg("redis command")
	return resp
}
