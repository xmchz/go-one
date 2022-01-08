package log

import (
	"os"
	"sync"

	"github.com/xmchz/go-one/log/core"
	"github.com/xmchz/go-one/log/formatter"
	"github.com/xmchz/go-one/log/writer"
)

var (
	once          sync.Once
	defaultLogger = core.New(
		core.WithWriters(writer.NewConsole(&formatter.Text{})),
		core.WithCallerSkip(2),
	)
)

func Init(opts ...core.Option) {
	once.Do(func() {
		defaultLogger = core.New(append(opts, core.WithCallerSkip(2))...)
	})
}

func Stop() {
	defaultLogger.Stop()
}

func Debug(format string, args ...interface{}) {
	defaultLogger.LogMsg(core.DebugLevel, format, args...)
}

func Trace(format string, args ...interface{}) {
	defaultLogger.LogMsg(core.TraceLevel, format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.LogMsg(core.InfoLevel, format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.LogMsg(core.WarnLevel, format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.LogMsg(core.ErrorLevel, format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.LogMsg(core.ErrorLevel, format, args...)
	defaultLogger.Stop()
	os.Exit(1)
}
