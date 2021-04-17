package breaker

import "github.com/xmchz/go-one/log"

type Acceptable func(err error) bool

type Breaker interface {
	Do(req func() error, acceptable Acceptable) error
}

type Empty struct {

}

func (b *Empty) Do(req func() error, acceptable Acceptable) error {
	log.Debug("empty breaker do with acceptable")
	return req()
}

