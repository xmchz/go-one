package access

import (
	"github.com/casbin/casbin/v2/persist"
)

type Storage interface {
	AccessAdapter(tblName string) (Adapter, error)
}

type Adapter interface {
	persist.Adapter
}
