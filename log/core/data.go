package core

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
	Fields   fields    `json:"fields,omitempty"`
}

type fields map[string]interface{}

func (fs fields) String() string {
	var buf bytes.Buffer
	for k, v := range fs {
		buf.WriteString(fmt.Sprintf("%v=%v ", k, v))
	}
	return buf.String()
}
