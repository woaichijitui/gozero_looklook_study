package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gozero_looklook_study/app/usercenter/cmd/api/internal/config"
	"gozero_looklook_study/app/usercenter/cmd/rpc/pb"
	"gozero_looklook_study/app/usercenter/cmd/rpc/usercenter"
)

type ServiceContext struct {
	Config        config.Config
	UserCenterRpc pb.UsercenterClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		UserCenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UserCenterRpc)),
	}
}
