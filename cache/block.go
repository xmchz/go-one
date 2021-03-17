package cache

import (
	"errors"
	"golang.org/x/sync/singleflight"
)

type Block struct {
	Cache
	singleflight.Group
}

func (c *Block) Take(dest interface{}, key string, query func(v interface{}) error) error {
	err := c.Get(key, dest)
	if errors.Is(ErrNotFound, err) {
		dest, err, _ = c.Do(key, func() (interface{}, error) {
			var v interface{}
			err = query(v)
			if err == nil {
				_ = c.Set(key, v)
			}
			return v, err
		})
	}
	return err
}
