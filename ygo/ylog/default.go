package ylog

import (
	"context"
	"fmt"
)

var zLogInstance ILogger = new(defaultLog)

type defaultLog struct{}

func (log *defaultLog) InfoF(format string, v ...interface{}) {
	StdYLog.Infof(format, v...)
}

func (log *defaultLog) ErrorF(format string, v ...interface{}) {
	StdYLog.Errorf(format, v...)
}

func (log *defaultLog) DebugF(format string, v ...interface{}) {
	StdYLog.Debugf(format, v...)
}

func (log *defaultLog) InfoFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	StdYLog.Infof(format, v...)
}

func (log *defaultLog) ErrorFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	StdYLog.Errorf(format, v...)
}

func (log *defaultLog) DebugFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	StdYLog.Debugf(format, v...)
}

func SetLogger(newlog ILogger) {
	zLogInstance = newlog
}

func Ins() ILogger {
	return zLogInstance
}
