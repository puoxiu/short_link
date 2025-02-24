package svc

import (
	"short_link_pro/sl_auth/auth_api/internal/config"
	"short_link_pro/sl_auth/auth_models"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModel auth_models.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.UserDB.DSN)

	return &ServiceContext{
		Config: c,
		UserModel: auth_models.NewUserModel(conn),
	}
}
