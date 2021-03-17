package writer

import (
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/xmchz/go-one/log"
	"io"
	"os"
	"path"
	"time"
)

func NewFile(logName, relativePath string, options ...Option) (*file, error) {
	fw := &file{
		logName:      logName,
		relativePath: relativePath,
	}
	for _, opt := range options {
		opt(fw)
	}
	if err := fw.init(); err!= nil {
		return nil, err
	}
	return fw, nil
}

type Option func(*file)

func WithRotateTime(rotateTime time.Duration) Option {
	return func(fw *file) {
		fw.rotateTime = rotateTime
	}
}

func WithMaxAge(maxAge time.Duration) Option {
	return func(fw *file) {
		fw.maxAge = maxAge
	}
}

func WithFormatter(formatter log.Formatter) Option {
	return func(fw *file) {
		fw.Formatter = formatter
	}
}

type file struct {
	io.WriteCloser
	log.Formatter
	logName      string
	relativePath string
	rotateTime   time.Duration
	maxAge       time.Duration
}

func (fw *file) Write(data *log.Data) {
	_, _ = fw.WriteCloser.Write(fw.Format(data))
}

func (fw *file) Close() {
	_ = fw.WriteCloser.Close()
}

func (fw *file) init() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	logPath := path.Join(pwd, fw.relativePath)
	if _, err := os.Stat(logPath); err != nil {
		if err := os.Mkdir(logPath, os.ModePerm); err != nil {
			return err
		}
	}
	rotate, err := rotateLogs.New(
		path.Join(logPath, fw.logName),
		rotateLogs.WithRotationTime(fw.rotateTime),
		rotateLogs.WithMaxAge(fw.maxAge),
	)
	if err != nil {
		return err
	}
	fw.WriteCloser = rotate
	return nil
}
