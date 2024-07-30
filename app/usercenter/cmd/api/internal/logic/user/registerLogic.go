package user

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gozero_looklook_study/app/usercenter/cmd/model"
	"gozero_looklook_study/app/usercenter/cmd/rpc/pb"

	"gozero_looklook_study/app/usercenter/cmd/api/internal/svc"
	"gozero_looklook_study/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// register
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	registerRpcResp, err := l.svcCtx.UserCenterRpc.Register(l.ctx, &pb.RegisterReq{
		Mobile:   req.Mobile,
		Password: req.Password,
		AuthKey:  req.Mobile,
		AuthType: model.UserAuthTypeSystem,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	resp = &types.RegisterResp{
		AccessToken:  registerRpcResp.AccessToken,
		AccessExpire: registerRpcResp.AccessExpire,
		RefreshAfter: registerRpcResp.RefreshAfter,
	}
	fmt.Printf("%+v", resp)
	return resp, nil
}
