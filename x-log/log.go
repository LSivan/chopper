package xlog

import (
	"go.uber.org/zap"
)

var sugar = func() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	return logger.Sugar()
}()

var Sugar = &Logger{sugar}

type Logger struct {
	*zap.SugaredLogger
}

func (s *Logger) Named(name string) Logger {
	return Logger{sugar.Named(name)}
}
