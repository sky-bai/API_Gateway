package serviceInfo

import (
	"API_Gateway/api/internal/global"
	"context"
	"time"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ServiceFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServiceFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) ServiceFlowLogic {
	return ServiceFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ServiceFlow 大盘流量统计 对应首页今日流量统计
func (l *ServiceFlowLogic) ServiceFlow() (*types.ServiceFlowResponse, error) {
	counter, err := global.FlowCounterHandler.GetCounter(global.FlowTotal)
	if err != nil {
		return nil, err
	}

	var todayList []int
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, time.Local)
		hourData, _ := counter.GetHourData(dateTime) // 不处理这里的错误 因为没有值就是0 但是可能就是redis报错
		todayList = append(todayList, int(hourData))
	}

	var yesterdayList []int
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, time.Local)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, int(hourData))
	}

	return &types.ServiceFlowResponse{
		Today:     todayList,
		Yesterday: yesterdayList}, nil
}
