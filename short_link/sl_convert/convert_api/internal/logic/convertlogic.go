package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"short_link_pro/pkg/base62"
	"short_link_pro/pkg/connect"
	md5 "short_link_pro/pkg/mds"
	"short_link_pro/pkg/urltool"
	"short_link_pro/shorturlmapmodel"
	"short_link_pro/sl_convert/convert_api/internal/svc"
	"short_link_pro/sl_convert/convert_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}


// Convert 短链接转换: 长链接-->短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 1. 校验输入合法性
	// 1.1 数据非空
	// 1.2 输入的长链接必须是正常访问
	if !connect.Get(req.LongUrl) {
		// err = ErrUrlInvalid
		return nil, errors.New("url invalid")
	}

	// 1.3 长链接是否已经存在(md5)--查db
	md5 := md5.Cal(req.LongUrl)
	u, err := l.svcCtx.ShortUrlMapModel.FindOneByMd5(l.ctx, sql.NullString{String: md5, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			// 已经存在
			return nil, fmt.Errorf("url already exists: %s", u.Surl.String)
		}
		logx.Errorw("find one by md5 failed", logx.LogField{Key: "err", Value: err})
		return nil, err
	}

	// 1.4 输入的不能是一个已经生成的短链接, eg: baidu.com/1dvwd23
	basePath, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw("get base path failed", logx.LogField{Key: "err", Value: err})
		return nil, err
	}
	u, err = l.svcCtx.ShortUrlMapModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			// 已经是短链接
			return nil, fmt.Errorf("url already exists: %s", u.Surl.String)
		}
		logx.Errorw("find one by surl failed", logx.LogField{Key: "err", Value: err})
		return nil, err
	}
	
	// 2. 根据辅助mysql主键 取号
	var short string
	for {
		seq, err := l.svcCtx.Sequence.NextNumber()
		if err != nil {
			logx.Errorw("next number failed", logx.LogField{Key: "err", Value: err})
			return nil, err
		}
		fmt.Println("seq: ", seq)

		// 3. 号码转短链接（62进制）
		// todo 安全性考虑--修改62进制字符串打乱
		// 避免特殊单词
		short = base62.Int2String(seq)
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; !ok {
			break
		}
	}

	// 4. 存储
	if _, err := l.svcCtx.ShortUrlMapModel.Insert(
		l.ctx,
		&shorturlmapmodel.ShortUrlMap{
			Lurl: sql.NullString{String: req.LongUrl, Valid: true},
			Surl: sql.NullString{String: short, Valid: true},
			Md5:  sql.NullString{String: md5, Valid: true},
		},
	); err != nil {
		logx.Errorw("insert failed", logx.LogField{Key: "err", Value: err})
		return nil, err
	}

	// 将生成的短链接存入布隆过滤器
	if err = l.svcCtx.BloomFilter.Add(l.ctx, short); err != nil {
		logx.Errorw("add to bloom filter failed", logx.LogField{Key: "err", Value: err})
		return nil, err
	}

	// 5. 返回	
	resp = &types.ConvertResponse{
		ShortUrl: l.svcCtx.Config.ShortDomain + short,
	}
	return resp, nil
}
