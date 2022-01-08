package core_test

import (
	"testing"
	"time"

	"github.com/xmchz/go-one/log/core"
	"github.com/xmchz/go-one/log/formatter"
	"github.com/xmchz/go-one/log/writer"
)

func TestConsoleLogger(t *testing.T) {
	lg := core.New(
		core.WithWriters(writer.NewConsole(&formatter.Json{})),
	)
	defer lg.Stop()
	lg.LogMsg(core.InfoLevel, "this is info log message")
}

func TestFileWriter(t *testing.T) {
	w, err := writer.NewFile(
		"test.%Y%m%d%H%m%S.log",
		"logs",
		writer.WithRotateTime(8*time.Second),
		writer.WithMaxAge(10*time.Second),
		writer.WithFormatter(&formatter.Text{}),
	)
	if err != nil {
		t.Fatal(err)
	}
	lg := core.New(
		core.WithWriters(
			w,
			writer.NewConsole(&formatter.Json{}),
		),
	)
	defer lg.Stop()
	lg.LogFields(core.InfoLevel, map[string]interface{}{"name": "chouchou", "age": 1})
	// for i := 0; i < 10; i++ {
	// 	lg.LogMsg(core.InfoLevel, "this is info log message")
	// 	lg.LogFields(core.InfoLevel, map[string]interface{}{"name": "chouchou", "age": 1})
	// 	time.Sleep(1 * time.Second)
	// }
}
