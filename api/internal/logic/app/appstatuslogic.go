package app

import (
	"context"
	"time"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AppStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) AppStatusLogic {
	return AppStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AppStatus 服务详情
func (l *AppStatusLogic) AppStatus(req types.AppStatusRequest) (*types.AppResponse, error) {
	var todayList []int
	var yesterdayList []int
	for i := 0; i < time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}
	for i := 0; i < 24; i++ {
		yesterdayList = append(yesterdayList, 0)
	}

	return &types.AppResponse{}, nil
}
