package logic

import (
	"context"
	"github.com/pkg/errors"
	"gozero_looklook_study/common/tool"
	"gozero_looklook_study/common/xerr"
	"time"

	"gozero_looklook_study/app/usercenter/cmd/rpc/internal/svc"
	"gozero_looklook_study/app/usercenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateTokenLogic) GenerateToken(in *pb.GenerateTokenReq) (*pb.GenerateTokenResp, error) {
	now := time.Now().Unix()

	expire := l.svcCtx.Config.JwtAuth.AccessExpire

	accesstoken, err := tool.GenerateTokenUsingRS256(in.UserId, now, expire)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TOKEN_GENERATE_ERROR), "getJwtToken err userId:%d , err:%v", in.UserId, err)
	}

	return &pb.GenerateTokenResp{
		AccessToken:  accesstoken,
		AccessExpire: now + expire,
		RefreshAfter: now + expire/2,
	}, nil
}
