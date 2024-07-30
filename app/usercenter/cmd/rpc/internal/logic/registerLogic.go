package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gozero_looklook_study/app/usercenter/cmd/model"
	"gozero_looklook_study/app/usercenter/cmd/rpc/internal/svc"
	"gozero_looklook_study/app/usercenter/cmd/rpc/pb"
	"gozero_looklook_study/common/tool"
	"gozero_looklook_study/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserAlreadyRegisterError = xerr.NewErrMsg("user has been registered")
var ErrGenerateTokenError = xerr.NewErrMsg("生成token失败")
var ErrUsernamePwdError = xerr.NewErrMsg("账号或密码不正确")

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	//查找
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "mobile:%s,err:%v", in.Mobile, err)
	}
	if user != nil {
		return nil, errors.Wrapf(ErrUserAlreadyRegisterError, "Register user exists mobile:%s,err:%v", in.Mobile, err)
	}

	//插入
	var userId int64
	if err := l.svcCtx.UserModel.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		//	插入user表
		user := new(model.User)
		user.Mobile = in.Mobile
		user.Nickname = in.Nickname
		if len(in.Nickname) == 0 {
			user.Nickname = tool.Krand(8, tool.KC_RAND_KIND_ALL)
		}
		fmt.Println(111)
		if len(in.Password) > 0 {
			user.Password = tool.Md5ByString(in.Password)
		}
		result, err := l.svcCtx.UserModel.Insert(l.ctx, user)
		fmt.Println(err)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user Insert err:%v,user:%+v", err, user)
		}
		fmt.Println(333)
		lastId, err := result.LastInsertId()
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user insertResult.LastInsertId err:%v,user:%+v", err, user)
		}
		userId = lastId
		fmt.Println(user)
		//	插入user_auth表
		userAuth := new(model.UserAuth)
		userAuth.UserId = lastId
		userAuth.AuthKey = in.AuthKey
		userAuth.AuthType = in.AuthType

		fmt.Println(userAuth)
		if _, err := l.svcCtx.UserAuthModel.Insert(ctx, userAuth); err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user_auth Insert err:%v,userAuth:%v", err, userAuth)
		}

		return nil

	}); err != nil {
		return nil, err
	}

	//生成token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	generateTokenResp, err := generateTokenLogic.GenerateToken(&pb.GenerateTokenReq{
		UserId: userId,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "GenerateToken userId : %d", userId)
	}

	//响应

	return &pb.RegisterResp{
		AccessToken:  generateTokenResp.AccessToken,
		AccessExpire: generateTokenResp.AccessExpire,
		RefreshAfter: generateTokenResp.RefreshAfter,
	}, nil
}
