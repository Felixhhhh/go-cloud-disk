package svc

import (
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/rest"
	"go-cloud-disk/core/internal/config"
	"go-cloud-disk/core/internal/middleware"
	"go-cloud-disk/core/modules"
	"xorm.io/xorm"
)

type ServiceContext struct {
	Config config.Config
	Engine *xorm.Engine
	RDB    *redis.Client
	Auth   rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Engine: modules.InitDB(c.Mysql.DataSource),
		RDB:    modules.InitRDB(c),
		Auth:   middleware.NewAuthMiddleware().Handle,
	}
}
