package writer

import (
	// "fmt"
	"github.com/xmchz/go-one/log/core"
	"os"
)

func NewConsole(f core.Formatter) *console{
	return &console{f}
}

type console struct {
	core.Formatter
}

func (w *console) Write(data *core.Data) {
	colored := cm[data.Level].Add(string(w.Format(data)))
	_, _ = os.Stdout.Write(append([]byte(colored), '\n'))
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func (w *console) Close() {
	// os.Stdout.Close()
}

