package log

import (
	"bytes"
	"fmt"
	"time"
)

type Data struct {
	Level    Level     `json:"level"`
	Time     time.Time `json:"time"`
	Filename string    `json:"filename"`
	LineNo   int       `json:"line_no"`
	Message  string    `json:"message,omitempty"`
	Fields   Fields    `json:"fields,omitempty"`
}

type Fields map[string]interface{}

func (fields Fields) String() string {
	var buf bytes.Buffer
	for k, v := range fields {
		buf.WriteString(fmt.Sprintf("%v=%v ", k, v))
	}
	return string(buf.Bytes())
}
