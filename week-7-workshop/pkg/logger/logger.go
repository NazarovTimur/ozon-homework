package logger

import (
	"context"
	"go.uber.org/zap"
	"sync"
)

var (
	globalLogger *Logger
	once         = sync.Once{}
)

type Logger struct {
	l *zap.SugaredLogger
}

func NewLogger(config zap.Config) *Logger {
	l, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	once.Do(func() {
		globalLogger = &Logger{l: l.Sugar()}
	})

	return globalLogger
}

type CtxKey string

const LoggerCtxKey CtxKey = "logger"

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, LoggerCtxKey, logger)
}

func Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if l, ok := ctx.Value(LoggerCtxKey).(*Logger); ok && l != nil {
		l.l.Infow(msg, keysAndValues)
		return
	}

	if globalLogger == nil {
		panic("global logger is nil")
	}

	globalLogger.l.Infow(msg, keysAndValues...)
}

func Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if l, ok := ctx.Value(LoggerCtxKey).(*Logger); ok && l != nil {
		l.l.Errorw(msg, keysAndValues)
		return
	}

	if globalLogger == nil {
		panic("global logger is nil")
	}

	globalLogger.l.Errorw(msg, keysAndValues...)
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}
