package serviceInfo

import (
	"API_Gateway/api/internal/global"
	"context"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/stores/sqlc"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PanelDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPanelDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) PanelDataLogic {
	return PanelDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PanelData 首页指标统计
func (l *PanelDataLogic) PanelData() (*types.PanelDataOutput, error) {

	// 1.获取服务数量
	serviceNum, err := l.svcCtx.GatewayServiceInfoModel.GetServiceNum()
	if err != nil && err != sqlc.ErrNotFound {
		return nil, errors.New("获取服务数量失败")
	}

	// 2.获取租户数量
	appCount, err := l.svcCtx.GatewayAppModel.GetAppCount()
	if err != nil && err != sqlc.ErrNotFound {
		return nil, errors.New("获取服务数量失败")
	}

	counter, _ := global.FlowCounterHandler.GetCounter(global.FlowTotal)
	//fmt.Println("-----------------",counter.QPS)

	return &types.PanelDataOutput{
		ServiceNum:      int64(serviceNum),
		AppNum:          int64(appCount),
		CurrentQPS:      counter.QPS,
		TodayRequestNum: counter.TotalCount}, nil
}
