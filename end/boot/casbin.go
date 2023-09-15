package boot

import (
	"InstantMessaging/app/global"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter"
)

func InitCasbin() {
	var err error
	// 初始化Casbin
	a := gormadapter.NewAdapterByDB(g.MysqlDB)
	g.Enforcer, err = casbin.NewEnforcer("./manifest/config/model.conf", "./manifest/config/policy.psv", a)
	if err != nil {
		fmt.Println(err)
	}

}
