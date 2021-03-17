package formatter

import (
	"encoding/json"
	"github.com/xmchz/go-one/log"
)

type Json struct {
}

func (w *Json) Format(data *log.Data) []byte {
	bs, _ := json.Marshal(data)
	return append(bs, []byte(sepLine)...)
}
