package logic

import (
	"context"

	"short_link_svc/api/internal/svc"
	"short_link_svc/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateURLLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateURLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateURLLogic {
	return &CreateURLLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateURL 创建短链接
func (l *CreateURLLogic) CreateURL(req *types.CreateURLRequest) (resp *types.CreateURLResponse, err error) {
	

	return
}
