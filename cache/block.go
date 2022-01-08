package cache

import (
	"encoding/json"
	"errors"

	"github.com/xmchz/go-one/log"
	"golang.org/x/sync/singleflight"
)

type Block struct {
	Cache
	singleflight.Group
}

func (c *Block) Take(dest interface{}, key string, query func(v interface{}) error) error {
	err := c.Get(key, dest)
	if errors.Is(err, ErrNotFound) {

		sharedMarshaledVal, err, shared := c.Do(key, func() (interface{}, error) {
			// single flight
			err = query(dest)
			log.Debug("%s do: %v", key, dest)
			if err != nil {
				return nil, err
			}
			_ = c.Set(key, dest)
			return json.Marshal(dest)  // marshal to share value
		})
		if err != nil {
			return err
		}
		// concurrency shared
		json.Unmarshal(sharedMarshaledVal.([]byte), dest)
		if shared {
			log.Debug("%s shared: %v", key, sharedMarshaledVal)
		}
	}
	return err
}
