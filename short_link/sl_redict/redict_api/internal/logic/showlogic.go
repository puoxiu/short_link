package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"short_link_pro/pkg/ip"
	"short_link_pro/sl_redict/redict_api/constants"
	"short_link_pro/sl_redict/redict_api/internal/svc"
	"short_link_pro/sl_redict/redict_api/internal/types"
	"short_link_pro/sl_redict/redict_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}


// Show 查看短链接, 将输入的短链接重定向到原始链接
func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// 查询布隆过滤器
	exists, err := l.svcCtx.BloomFilter.Contains(l.ctx, req.ShortUrl)
	if err != nil {
		logx.Errorw("查询布隆过滤器失败", logx.LogField{Key: "err", Value: err})
		return nil, err
	}
	if !exists {
		return nil, errors.New("404, 短链接不存在1")
	}

	// 根据短链接查询原始链接
	u, err := l.svcCtx.ShortUrlMapModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("404, 短链接不存在2")
		}
		logx.Errorw("查询短链接失败", logx.LogField{Key: "err", Value: err})
		return nil, err
	}
	clientIp, clientAgent := "unknown", "unknown"
	if ip, ok := l.ctx.Value(constants.UserIPKey).(string); ok {
		clientIp = ip
	}
	if agent, ok := l.ctx.Value(constants.UserAgentKey).(string); ok {
		clientAgent = agent
	}
	go saveAccessLog2DB(l.svcCtx, u.Id, clientIp, clientAgent)
	
	resp = &types.ShowResponse{
		LongUrl: u.Lurl.String,
	}
	return resp, nil
}

// saveAccessLog2DB 保存访问日志到数据库
func saveAccessLog2DB(svcCtx *svc.ServiceContext, shortUrlId int64, clientIp, clientAgent string) {
	country, region, city, err := ip.GetLocation(clientIp)
	if err != nil {
		logx.Errorw("获取访问者信息失败", logx.LogField{Key: "err", Value: err})
		return
	}
	fmt.Println("ShortUrlId:", shortUrlId)
	fmt.Println("ip:", clientIp)
	fmt.Println("agent:", clientAgent)
	fmt.Println("country:", country)
	fmt.Println("region:", region)
	fmt.Println("city:", city)

	// 存储
	dbCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	fmt.Println("正常")
	if _, err := svcCtx.ShortUrlAccessLogModel.Insert(
		dbCtx, 
		&redict_models.ShortUrlAccessLog{
		ShortUrlId: shortUrlId,
		UserAgent: sql.NullString{String: clientAgent, Valid: true},
		AccessIp: clientIp,
		AccessTime: time.Now(),
		Country: sql.NullString{String: country, Valid: true},
		Region: sql.NullString{String: region, Valid: true},
		City: sql.NullString{String: city, Valid: true},
	}); err != nil {
		logx.Error("保存访问日志失败", logx.LogField{Key: "err", Value: err})
		return
	}
}