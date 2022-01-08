package writer

import (
	"fmt"
	"github.com/xmchz/go-one/log/core"
)

type color uint8

const (
	black color = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

func (c color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

var cm = map[core.Level]color{
	core.InfoLevel:   green,
	core.DebugLevel:  cyan,
	core.TraceLevel:  white,
	core.WarnLevel:   yellow,
	core.ErrorLevel:  red,
}
