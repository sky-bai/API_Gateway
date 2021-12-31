package app

import (
	"API_Gateway/api/internal/global"
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
func (l *AppStatusLogic) AppStatus(req types.AppStatusRequest) (*types.AppStatus, error) {

	// 1.通过Id获取app_id
	appInfo, err := l.svcCtx.GatewayAppModel.FindOne(req.ID)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	counter, err := global.FlowCounterHandler.GetCounter(FlowAppPrefix + appInfo.AppId)
	if err != nil {
		return nil, err
	}

	var todayList []int
	currentTime := time.Now()
	for i := 0; i < time.Now().Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, time.Local)
		hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, int(hourData))
	}

	var yesterdayList []int
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i < 24; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, time.Local)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, int(hourData))
	}

	return &types.AppStatus{
		Today:     todayList,
		Yesterday: yesterdayList}, nil
}
