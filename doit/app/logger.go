package app

import (
	"io"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/rs/zerolog"
)

func NewLeveledLogger(path string) *LeveledLogger {
	l := new(LeveledLogger)
	l.Writer = os.Stderr
	l.targets = make(map[string]*lumberjack.Logger)
	for _, level := range []zerolog.Level{
		zerolog.DebugLevel,
		zerolog.InfoLevel,
		zerolog.WarnLevel,
		zerolog.ErrorLevel,
		zerolog.FatalLevel,
		zerolog.PanicLevel,
		zerolog.NoLevel,
	} {
		name := l.parseLevel(level)
		l.targets[name] = &lumberjack.Logger{
			Filename:   path + "/" + name + ".log",
			MaxSize:    100, // 兆字节
			MaxBackups: 10,  //	保留的最大旧日志文件数
			MaxAge:     28,  //	最大保留时间天数
		}
	}
	return l
}

type LeveledLogger struct {
	io.Writer
	targets map[string]*lumberjack.Logger
}

func (l *LeveledLogger) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	return l.targets[l.parseLevel(level)].Write(p)
}

func (l LeveledLogger) parseLevel(level zerolog.Level) string {
	if name := level.String(); name != "" {
		return name
	}
	return "default"
}
