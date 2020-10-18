package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"os"
	"strings"
	"time"
)

type Level int8

const (
	TraceLevel Level = iota - 2
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

type Logger interface {
	Log(level Level, v ...interface{})
	Logf(level Level, format string, v ...interface{})
	LogFields(fields map[string]interface{})
}

func New(w io.Writer) Logger {
	return newLogrus(map[Level]io.Writer{
		InfoLevel: w,
		ErrorLevel: w,
	})
}

func NewRotateLog(logName, relativePath string) Logger {
	w, _ := newRotateWriter(logName, relativePath)
	return newLogrus(map[Level]io.Writer{
		InfoLevel: w,
		ErrorLevel: w,
	})
}

/*
logName 日志名
logPath 相对路径
*/
func newRotateWriter(logName, relativePath string) (io.Writer, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	logPath := strings.Join([]string{pwd, relativePath}, string(os.PathSeparator))
	if _, err := os.Stat(logPath); err != nil {
		if err := os.Mkdir(logPath, os.ModePerm); err != nil {
			return nil, err
		}
	}
	return rotatelogs.New(
		logPath+string(os.PathSeparator)+logName+".%Y%m%d.log",
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)
}