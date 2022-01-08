package formatter

import (
	"encoding/json"
	"github.com/xmchz/go-one/log/core"
)

type Json struct {
}

func (w *Json) Format(data *core.Data) []byte {
	bs, _ := json.Marshal(data)
	return bs
}
