package log_test

import (
	"github.com/xmchz/go-one/log"
	"github.com/xmchz/go-one/log/formatter"
	"github.com/xmchz/go-one/log/writer"
	"testing"
	"time"
)

func TestConsoleLogger(t *testing.T) {
	log.Init(
		log.WithWriters(&writer.Console{Formatter: &formatter.Text{}}),
	)
	log.Info("this is info log message")
	log.Error("this is error log message")
	log.Warn("this is warn log message")
	time.Sleep(1 * time.Second)
}

func TestFileWriter(t *testing.T) {
	w, err := writer.NewFile("test.%Y%m%d%H%m%S",
		"logs",
		writer.WithRotateTime(10*time.Second),
		writer.WithMaxAge(15*time.Second),
		writer.WithFormatter(&formatter.Text{}),
	)
	if err != nil {
		t.Fatal(err)
	}
	log.Init(
		log.WithWriters(w, &writer.Console{Formatter: &formatter.Text{}}),
	)
	for i := 0; i < 30; i++ {
		log.Info("this is info log message")
		time.Sleep(1 *time.Second)
	}
}

func TestNewLogger(t *testing.T) {
	w, err := writer.NewFile("access.%Y%m%d",
		"logs",
		writer.WithFormatter(&formatter.Text{}),
	)
	if err != nil {
		t.Fatal(err)
	}
	lg := log.New(
		log.WithWriters(w, &writer.Console{Formatter: &formatter.Json{}}),
	)
	lg.Log(log.InfoLevel, "this is log info")
	lg.LogFields(map[string]interface{}{"name": "chouchou", "age": 1})
	time.Sleep(1 *time.Second)
}
