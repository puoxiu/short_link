package logic

import (
	"context"
	"errors"

	jwts "short_link_pro/pkg/jwt"
	md5 "short_link_pro/pkg/mds"
	"short_link_pro/sl_auth/auth_api/internal/svc"
	"short_link_pro/sl_auth/auth_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)


var (
	ErrUserNotExist = errors.New("用户不存在")
	ErrPassword	 = errors.New("密码错误")
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	if req.Username == "" || req.Password == "" {
		return nil, ErrParam
	}

	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		if err == sqlx.ErrNotFound {
			logx.Errorw("find one by username failed", logx.LogField{Key: "err", Value: err})
			return nil, ErrUserNotExist
		}
		return nil, ErrServer
	}
	if u.Password != md5.EncryptPassword(req.Password) {
		return nil, ErrPassword
	}

	token,err := jwts.GenToken(jwts.JwtPayLoad{
		UserID: uint(u.Id),
		Username: u.Username,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Errorw("gen token failed", logx.LogField{Key: "err", Value: err})
		return nil, ErrServer
	}

	resp = &types.LoginResponse{
		Token: token,
	}
	return resp, nil
}
