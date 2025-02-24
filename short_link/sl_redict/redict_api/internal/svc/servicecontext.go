package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"short_link_pro/sl_redict/redict_api/internal/config"
	"short_link_pro/sl_redict/redict_api/internal/middleware"
)

type ServiceContext struct {
	Config             config.Config
	ClientIPMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		ClientIPMiddleware: middleware.NewClientIPMiddleware().Handle,
	}
}
