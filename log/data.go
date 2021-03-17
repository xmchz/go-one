package log

import (
	"bytes"
	"fmt"
	"time"
)

type Data struct {
	Level    Level     `json:"level"`
	Time     time.Time `json:"time"`
	Message  string    `json:"message"`
	Filename string    `json:"filename"`
	LineNo   int       `json:"line_no"`
	Fields   Fields
}

type Fields map[string]interface{}

func (fields Fields) String() string {
	var buf bytes.Buffer
	for k, v := range fields {
		buf.WriteString(fmt.Sprintf("%v=%v ", k, v))
	}
	return string(buf.Bytes())
}
