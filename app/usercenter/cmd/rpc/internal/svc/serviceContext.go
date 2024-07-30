package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gozero_looklook_study/app/usercenter/cmd/model"
	"gozero_looklook_study/app/usercenter/cmd/rpc/internal/config"
)

type ServiceContext struct {
	Config        config.Config
	UserModel     model.UserModel
	UserAuthModel model.UserAuthModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		UserModel:     model.NewUserModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		UserAuthModel: model.NewUserAuthModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
	}
}
