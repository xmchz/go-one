package formatter

import (
	"bytes"
	"fmt"
	"github.com/xmchz/go-one/log"
)

const (
	sepSpace   = " "
	sepColon   = ":"
	sepLine    = "\n"
	layoutTime = "2006-01-02 15:04:05.999"
)

type Text struct {
}

func (w *Text) Format(data *log.Data) []byte {
	var buf bytes.Buffer
	buf.WriteString(data.Time.Format(layoutTime))
	buf.WriteString(sepSpace)
	buf.WriteString(data.Level.String())
	buf.WriteString(sepSpace)
	buf.WriteString(fmt.Sprintf("%s%s%d", data.Filename, sepColon, data.LineNo))
	buf.WriteString(sepSpace)
	buf.WriteString(data.Fields.String())
	buf.WriteString(sepSpace)
	buf.WriteString(data.Message)
	buf.WriteString(sepLine)
	return buf.Bytes()
}
