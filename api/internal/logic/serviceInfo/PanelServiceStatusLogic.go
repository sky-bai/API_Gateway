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

type PanelServiceStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPanelServiceStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) PanelServiceStatusLogic {
	return PanelServiceStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PanelServiceStatus 大盘服务状态
func (l *PanelServiceStatusLogic) PanelServiceStatus() (*types.PanelServiceStatusResponse, error) {

	// 1.获取每个服务对应的数量
	allServiceNum, err := l.svcCtx.GatewayServiceInfoModel.GetAllNum()
	if err != nil && err != sqlc.ErrNotFound {
		return nil, err
	}

	var Legend []string

	for _, serviceNum := range allServiceNum {
		//fmt.Println(serviceNum.LoadType)
		serviceName, ok := global.LoadTypeMap[serviceNum.LoadType]
		if !ok {
			return nil, errors.New("未找到该服务")
		}
		//allServiceNum[index].Name = serviceName
		Legend = append(Legend, serviceName)
	}

	var data []types.DashServiceStatItem
	for _, service := range allServiceNum {
		data = append(data, types.DashServiceStatItem{
			Name:     global.LoadTypeMap[service.LoadType],
			Value:    service.Value,
			LoadType: service.LoadType,
		})
	}

	return &types.PanelServiceStatusResponse{Legend: Legend, Data: data}, nil
}
