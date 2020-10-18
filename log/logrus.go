package log

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
)

var levelToLogrus = map[Level]logrus.Level{
	TraceLevel: logrus.TraceLevel,
	DebugLevel: logrus.DebugLevel,
	InfoLevel:  logrus.InfoLevel,
	WarnLevel:  logrus.WarnLevel,
	ErrorLevel: logrus.ErrorLevel,
	PanicLevel: logrus.PanicLevel,
	FatalLevel: logrus.FatalLevel,
}

func newLogrus(wm map[Level]io.Writer) *logrusLogger {
	lwm := make(map[logrus.Level]io.Writer)
	for l, w := range wm {
		lwm[levelToLogrus[l]] = w
	}
	lg := logrus.New()
	lg.AddHook(lfshook.NewHook(
		lfshook.WriterMap(lwm),
		&logrus.TextFormatter{},
	))
	return &logrusLogger{lg}
}

func newLogrusDefault() *logrusLogger {
	return &logrusLogger{logrus.StandardLogger()}
}

type logrusLogger struct {
	*logrus.Logger
}

func (l *logrusLogger) Log(level Level, v ...interface{}) {
	l.Logger.Log(levelToLogrus[level], v...)
}

func (l *logrusLogger) Logf(level Level, format string, v ...interface{}) {
	l.Logger.Logf(levelToLogrus[level], format, v...)
}

func (l *logrusLogger) LogFields(fields map[string]interface{}) {
	l.Logger.WithFields(fields).Print()
}
