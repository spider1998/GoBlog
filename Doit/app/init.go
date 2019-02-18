package app

import (
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/rs/zerolog"
	"os"
	"Project/Doit/resource"
)

var (
	Conf   Config         // 系统配置
	Logger zerolog.Logger // 全局日志
	DB     *dbx.DB        // 全局 DB 实例
	Redis  RedisBox       // 全局 redis 实例
)

//-----------------------------初始化配置----------------------------------------------------------------------

func Init() error {
	var err error

	Conf, err = LoadConfig()
	if err != nil {
		return err
	}

	//-----------------------------加载BOX(resource下的相关数据源文件)----------------------------------------------------------------------

	{
		resource.Load()
	}

	//-----------------------------配置日志及日志文件存储----------------------------------------------------------------------
	leveledLogger := NewLeveledLogger(Conf.ConfPath + "/logs")
	if Conf.Debug {
		Logger = zerolog.New(zerolog.MultiLevelWriter(leveledLogger, zerolog.ConsoleWriter{Out: os.Stderr})).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	} else {
		Logger = zerolog.New(leveledLogger).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	}
	//调用进程ID
	Logger.Info().Int("pid", os.Getpid()).Msg("app booted.")

	Logger.Info().Interface("config", Conf).Msg("loaded config.")

	//-----------------------------连接数据库----------------------------------------------------------------------

	Logger.Info().Msg("start load db...")
	DB, err = LoadDB(Conf.Mysql)
	if err != nil {
		return err
	}
	Logger.Info().Msg("loaded db.")

	//-----------------------------映射数据表----------------------------------------------------------------------

	Logger.Info().Msg("migrate db...")
	err = Migrate(Conf.Mysql)
	if err != nil {
		return err
	}
	Logger.Info().Msgf("applied migrations...")

	//-----------------------------连接并创建redis----------------------------------------------------------------------

	Logger.Info().Msg("try to load redis.")
	Redis, err = LoadRedis(Conf.Redis, Logger)
	if err != nil {
		return err
	}
	Logger.Info().Msg("loaded redis.")


	return nil
}
