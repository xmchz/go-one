package core

import (
	"fmt"
	"path"
	"runtime"
	"sync"
	"time"
)

const chSize = 1024

type Logger interface {
	LogMsg(level Level, format string, args ...interface{})
	LogFields(level Level, fields map[string]interface{})
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
		dataCh:     make(chan *Data, chSize),
		callerSkip: 1,
		level:      InfoLevel,
		middleCh:   make(chan *Data),
		closing:    make(chan struct{}),
		closed:     make(chan struct{}),
	}
	for _, opt := range opts {
		opt(lg)
	}
	lg.wg.Add(1)
	go lg.relay()
	go lg.write()
	return lg
}

type logger struct {
	wg         sync.WaitGroup
	dataCh     chan *Data
	level      Level
	writers    []Writer
	callerSkip int

	middleCh chan *Data
	closing  chan struct{}
	closed   chan struct{}
}

func (l *logger) LogMsg(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	_, fileName, lineNo, _ := runtime.Caller(l.callerSkip)
	l.send(&Data{
		Level:    level,
		Time:     time.Now(),
		Filename: path.Base(fileName),
		LineNo:   lineNo,
		Message:  fmt.Sprintf(format, args...),
	})
}

func (l *logger) LogFields(level Level, fields map[string]interface{}) {
	if level < l.level {
		return
	}
	_, fileName, lineNo, _ := runtime.Caller(l.callerSkip)
	l.send(&Data{
		Level:    level,
		Time:     time.Now(),
		Filename: path.Base(fileName),
		LineNo:   lineNo,
		Fields:   fields,
	})
}

// 停止工作并等待数据write完成
func (l *logger) Stop() {
	select {
	case l.closing <- struct{}{}: // 当前G通知关闭，等待关闭成功
		<-l.closed
	case <-l.closed: // 其他G通知过了，等待
	}
	l.wg.Wait()
	for _, w := range l.writers {
		w.Close()
	}
}

func (l *logger) send(v *Data) {
	select {
	case <-l.closed:
		return
	default:
	}
	select {
	case <-l.closed:
		return
	case l.middleCh <- v:
	}
}

func (l *logger) relay() {
	for {
		select {
		case <-l.closing: // 没有数据时，收到停止通知（没有该case可能阻塞在middleCh读取，因为middleCh不会关闭）
			close(l.closed)
			close(l.dataCh)
			return
		case v := <-l.middleCh:
			select {
			case <-l.closing: // 有待转发数据时，收到停止通知
				close(l.closed)
				l.dataCh <- v
				close(l.dataCh)
				return
			case l.dataCh <- v:
			}
		}
	}
}

func (l *logger) write() {
	// DO NOT add wg here
	defer l.wg.Done()
	for data := range l.dataCh {
		for _, w := range l.writers {
			w.Write(data)
		}
	}
}
