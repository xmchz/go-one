package rbac_test

import (
	"testing"

	"github.com/xmchz/go-one/access/rbac"
	"github.com/xmchz/go-one/log"
	"github.com/xmchz/go-one/storage"
)

type dbConf struct{}

func (d dbConf) MigrationUrl() string {
	return ""
}
func (d dbConf) MigrationVersion() uint {
	return 0
}

func (d dbConf) DbDriver() string {
	return "mysql"
}

func (d dbConf) DbName() string {
	return "chz_wiki"
}

func (d dbConf) DbUrl() string {
	return "root:chouchou@tcp(127.0.0.1:3306)/chz_wiki?charset=utf8mb4&parseTime=True"
}

func (d dbConf) DbInfo() string {
	return "root@127.0.0.1:3306/chz_wiki"
}

func (d dbConf) MaxIdleConn() int {
	return 5
}

func (d dbConf) MaxOpenConn() int {
	return 5
}

type accessConfFunc func() string
func(c accessConfFunc) AccessInfoTblName() string {
	return c()
}

func TestAccess(t *testing.T) {
	defer log.Stop()
	access := rbac.New(accessConfFunc(func() string{
		return "access_rule"
	}), storage.New(&dbConf{}))
	log.Info("%#v", access)
}