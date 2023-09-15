package g

import (
	"InstantMessaging/app/internal/model/config"
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

var (
	Config   *config.Config
	Logger   *zap.SugaredLogger
	MysqlDB  *gorm.DB
	Rdb      *redis.Client
	Enforcer *casbin.Enforcer
)
