package rbac

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/xmchz/go-one/access"
	"github.com/xmchz/go-one/log"
)

const (
	modeConf = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && regexMatch(r.act, p.act)
`
)

func New(conf access.Config, storage access.Storage) *Access {
	rbacModel, err := model.NewModelFromString(modeConf)
	if err != nil {
		log.Fatal("new model: %s", err.Error())
	}

	adapter, err := storage.AccessAdapter(conf.AccessInfoTblName())
	if err != nil {
		log.Fatal("new adapter: %s", err.Error())
	}
	enforcer, err := casbin.NewEnforcer(rbacModel, adapter)
	if err != nil {
		log.Fatal("new enforcer: %s", err.Error())
	}
	log.Info("rbac init success, access info table: %s", conf.AccessInfoTblName())
	return &Access{
		enforcer,
	}
}

type Access struct {
	*casbin.Enforcer
}
