package formatter

import (
	"bytes"
	"fmt"
	"github.com/xmchz/go-one/log/core"
)

const (
	sepSpace   = ' '
	sepColon   = ':'
	sepLine    = '\n'
	layoutTime = "2006-01-02 15:04:05.999"
)

type Text struct {
}

func (w *Text) Format(data *core.Data) []byte {
	var buf bytes.Buffer
	buf.WriteString(data.Time.Format(layoutTime))
	buf.WriteByte(sepSpace)
	buf.WriteString(data.Level.String())
	buf.WriteByte(sepSpace)
	buf.WriteString(fmt.Sprintf("%s%c%d", data.Filename, sepColon, data.LineNo))
	buf.WriteByte(sepSpace)
	buf.WriteString(data.Fields.String())
	buf.WriteByte(sepSpace)
	buf.WriteString(data.Message)
	return buf.Bytes()
}
