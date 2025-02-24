package logic

import (
	"context"
	"errors"

	md5 "short_link_pro/pkg/mds"
	"short_link_pro/sl_auth/auth_api/internal/svc"
	"short_link_pro/sl_auth/auth_api/internal/types"
	"short_link_pro/sl_auth/auth_models"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)
var (
	ErrParam = errors.New("param error")
	ErrNot = errors.New("两次密码不一致")
	ErrServer = errors.New("服务器内部错误")
	ErrorUserExist = errors.New("用户名已存在")
)
type SignupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}


func (l *SignupLogic) Signup(req *types.SignupRequest) (resp *types.SignupResponse, err error) {
	// 参数校验 todo 正则表达式校验
	if req.Username == "" || req.Password == "" {
		return nil, ErrParam
	}
	if req.Password != req.RePassword {
		return nil, ErrNot
	}

	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil && err != sqlx.ErrNotFound {
		logx.Errorw("find one by username failed", logx.LogField{Key: "err", Value: err})
		return nil, ErrServer
	}
	if err == nil || u != nil{
		return nil, ErrorUserExist
	}

	pwd := md5.EncryptPassword(req.Password)

	if _, err = l.svcCtx.UserModel.Insert(l.ctx, &auth_models.User{
		Username: req.Username,
		Password: pwd,
	}); err != nil {
		logx.Errorw("insert user failed", logx.LogField{Key: "err", Value: err})
		return nil, ErrServer
	}	
	
	resp = &types.SignupResponse{
		Message: "注册成功",
	}
	return resp, nil
}