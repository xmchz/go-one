package log

import (
	"fmt"
	"path"
	"runtime"
	"sync"
	"time"
)

type Logger interface {
	Log(level Level, format string, args ...interface{})
	LogFields(fields map[string]interface{})
}

type Writer interface {
	Write(*Data)
	Close()
}

type Formatter interface {
	Format(*Data) []byte
}

func New(opts ...Option) *logger {
	lg := &logger{
		dataCh:     make(chan *Data, defaultLoggerCh),
		callerSkip: 1,
		level: InfoLevel,
	}
	for _, opt := range opts {
		opt(lg)
	}
	lg.wg.Add(1)
	go lg.run()
	return lg
}

type logger struct {
	level      Level
	writers    []Writer
	dataCh     chan *Data
	callerSkip int
	wg         sync.WaitGroup
}

func (l *logger) Log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	_, fileName, lineNo, _ := runtime.Caller(l.callerSkip)
	data := &Data{
		Level:    level,
		Time:     time.Now(),
		Filename: path.Base(fileName),
		LineNo:   lineNo,
		Message:  fmt.Sprintf(format, args...),
	}
	select {
	case l.dataCh <- data:
	default:
		return
	}
}

func (l *logger) LogFields(fields map[string]interface{}) {
	_, fileName, lineNo, _ := runtime.Caller(l.callerSkip)
	data := &Data{
		Level:    InfoLevel,
		Time:     time.Now(),
		Filename: path.Base(fileName),
		LineNo:   lineNo,
		Fields:   fields,
	}
	select {
	case l.dataCh <- data:
	default:
		return
	}
}

func (l *logger) run() {
	for data := range l.dataCh {
		for _, w := range l.writers {
			w.Write(data)
		}
	}
	l.wg.Done()
}
