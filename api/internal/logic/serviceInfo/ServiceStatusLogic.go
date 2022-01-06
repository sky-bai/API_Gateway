package serviceInfo

import (
	"API_Gateway/api/internal/global"
	"context"
	"time"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

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

// ServiceStatus 获取单个服务流量统计
func (l *ServiceStatusLogic) ServiceStatus(req types.ServiceStatusResquest) (*types.ServiceStatusResponse, error) {

	// 1. 获取服务名字
	serviceInfo, err := l.svcCtx.GatewayServiceInfoModel.FindOne(int64(req.ID))
	if err != nil {
		l.Error(err)
		return nil, err
	}

	counter, err := global.FlowCounterHandler.GetCounter(global.FlowServicePrefix + serviceInfo.ServiceName)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	var todayList []int
	currentTime := time.Now()
	for i := 0; i < currentTime.Hour(); i++ {
		newTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, time.Local)
		//fmt.Println("!!!!!!", newTime)
		data, _ := counter.GetHourData(newTime)

		todayList = append(todayList, int(data))
	}

	var yesterdayList []int
	yesterTime := currentTime.Add(-24 * time.Hour)
	for i := 0; i < 24; i++ {
		newTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, time.Local)
		//fmt.Println("!!!!!!+++++++", newTime)
		data, _ := counter.GetHourData(newTime)

		yesterdayList = append(yesterdayList, int(data))
	}
	return &types.ServiceStatusResponse{Today: todayList, Yesterday: yesterdayList}, nil
}
