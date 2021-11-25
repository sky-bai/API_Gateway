package serviceInfo

import (
	"context"
	"time"

	"API_Gateway/api/interna/svc"
	"API_Gateway/api/interna/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ServiceStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServiceStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) ServiceStatusLogic {
	return ServiceStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 获取服务流量统计
func (l *ServiceStatusLogic) ServiceStatus(req types.ServiceDetailResquest) (*types.ServiceStatusResponse, error) {

	todayList := []int{}
	yesterdayList := []int{}
	for i := 0; i < time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}
	for i := 0; i < 24; i++ {
		yesterdayList = append(yesterdayList, 0)
	}
	return &types.ServiceStatusResponse{Today: todayList, Yesterday: yesterdayList}, nil
}
