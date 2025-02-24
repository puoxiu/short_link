package logic

import (
	"context"
	"errors"

	"short_link_pro/pkg"
	jwts "short_link_pro/pkg/jwt"
	"short_link_pro/sl_auth/auth_api/internal/svc"
	"short_link_pro/sl_auth/auth_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthenticationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	return &AuthenticationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthenticationLogic) Authentication(req *types.AuthenticationRequest) (resp *types.AuthenticationResponse, err error) {
	if pkg.InlistByRegs(l.svcCtx.Config.WhiteList, req.ValidPath) {
		logx.Infof("白名单请求：%s", req.ValidPath)
		return
	}

	if req.Token == "" {
		err = errors.New("认证失败1")
		return
	}
	claims, err := jwts.ParseToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		err = errors.New("认证失败2")
		return
	}

	resp = &types.AuthenticationResponse{
		Username: claims.Username,
	}

	return resp, nil
}

