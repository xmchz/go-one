package breaker

import "github.com/xmchz/go-one/log"

type Acceptable func(err error) bool

type Breaker interface {
	DoWithAcceptable(req func() error, acceptable Acceptable) error
}

type Empty struct {

}

func (b *Empty) DoWithAcceptable(req func() error, acceptable Acceptable) error {
	log.Debug("empty breaker do with acceptable")
	return req()
}

