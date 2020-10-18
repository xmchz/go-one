package log

import (
	"io"
	"os"
)

var defaultLog Logger = newLogrusDefault()

type Config interface {
	GetLogPath() string
	GetLogName() string
}

func Init(conf Config) error {
	logPath := conf.GetLogPath()
	logName := conf.GetLogName()
	errWriter, err := newRotateWriter(logName+".error", logPath)
	if err != nil {
		return err
	}
	infoWriter, err := newRotateWriter(logName+".info", logPath)
	if err != nil {
		return err
	}

	defaultLog = newLogrus(map[Level]io.Writer{
		InfoLevel:  infoWriter,
		ErrorLevel: errWriter,
	})
	Infof("%s default log init success", logName)
	return nil
}

func Debug(v ...interface{}) {
	defaultLog.Log(DebugLevel, v...)
}
func Debugf(format string, v ...interface{}) {
	defaultLog.Logf(DebugLevel, format, v...)
}
func Info(v ...interface{}) {
	defaultLog.Log(InfoLevel, v...)
}
func Infof(format string, v ...interface{}) {
	defaultLog.Logf(InfoLevel, format, v...)
}
func Warn(v ...interface{}) {
	defaultLog.Log(WarnLevel, v...)
}
func Warnf(format string, v ...interface{}) {
	defaultLog.Logf(WarnLevel, format, v...)
}
func Error(v ...interface{}) {
	defaultLog.Log(ErrorLevel, v...)
}
func Errorf(format string, v ...interface{}) {
	defaultLog.Logf(ErrorLevel, format, v...)
}
func Panic(v ...interface{}) {
	defaultLog.Log(PanicLevel, v...)
}
func Panicf(format string, v ...interface{}) {
	defaultLog.Logf(PanicLevel, format, v...)
}
func Fatal(v ...interface{}) {
	defaultLog.Log(FatalLevel, v...)
	os.Exit(1)
}
func Fatalf(format string, v ...interface{}) {
	defaultLog.Logf(FatalLevel, format, v...)
	os.Exit(1)
}
