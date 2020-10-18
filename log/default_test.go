package log

import (
	"testing"
)

func TestInfof(t *testing.T) {
	name := "test-log"
	Infof("%s log success", name)
}
