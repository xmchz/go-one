package writer

import (
	"fmt"
	"github.com/xmchz/go-one/log"
	"os"
)

type Console struct {
	log.Formatter
}

func (w *Console) Write(data *log.Data) {
	colored := cm[data.Level].Add(string(w.Format(data)))
	_, err := os.Stdout.Write([]byte(colored))
	if err != nil {
		fmt.Println(err)
	}
}

func (w *Console) Close() {
}

