package log

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	lg := New(os.Stdout)
	lg.LogFields(map[string]interface{}{
		"name": "log-test",
	})
}