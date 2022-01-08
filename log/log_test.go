package log_test

import (
	"testing"
	"time"

	"github.com/xmchz/go-one/log"
	"github.com/xmchz/go-one/log/core"
	"github.com/xmchz/go-one/log/formatter"
	"github.com/xmchz/go-one/log/writer"
)

func TestLog(t *testing.T) {
	defer log.Stop()
	log.Info("this is info message")
	// log.Warn("this is warn message")
	// log.Error("this is error message")
}

func TestLogFatal(t *testing.T) {
	log.Init(
		core.WithWriters(writer.NewConsole(&formatter.Json{})),
	)
	defer log.Stop()
	log.Info("this is info message")
	log.Fatal("this is fatal message") // how to test os.Exit(1)
}

func TestCustomLog(t *testing.T) {
	fileWriter, err := writer.NewFile(
		"test.%Y%m%d%H%m%S.log",
		"logs",
		writer.WithRotateTime(10*time.Second),
		writer.WithMaxAge(15*time.Second),
		writer.WithFormatter(&formatter.Text{}),
	)
	if err != nil {
		t.Fatal(err)
	}
	log.Init(
		core.WithWriters(
			fileWriter,
			writer.NewConsole(&formatter.Text{}),
		),
	)
	defer log.Stop()
	log.Info("this is info message")
}
