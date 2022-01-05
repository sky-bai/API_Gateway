package app

import (
	"API_Gateway/api/internal/global"
	"API_Gateway/model/ga_gateway_app"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/stores/sqlc"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AppListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppListLogic(ctx context.Context, svcCtx *svc.ServiceContext) AppListLogic {
	return AppListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const (
	FlowTotal         = "flow_total"
	FlowServicePrefix = "flow_service_"
	FlowAppPrefix     = "flow_app_"
)

func (l *AppListLogic) AppList(req types.AppListRequest) (*types.APPListResponse, error) {
	appInfo, err := l.svcCtx.GatewayAppModel.GetServiceList(req.Info, req.PageNo, req.PageSize)
	if err != nil && err != sqlc.ErrNotFound {
		fmt.Println(err)
		return nil, errors.New("获取租户列表失败")
	}

	appList := appInfo.Data.([]ga_gateway_app.GatewayApp)

	var list []types.APPListItemOutput
	for _, app := range appList {
		appCounter, err := global.FlowCounterHandler.GetCounter(FlowAppPrefix + app.AppId)
		if err != nil {
			return nil, errors.New("获取租户列表失败")
		}
		list = append(list, types.APPListItemOutput{
			ID:       app.Id,
			AppID:    app.AppId,
			Name:     app.Name,
			Secret:   app.Secret,
			WhiteIPS: app.WhiteIps,
			Qpd:      app.Qpd,
			Qps:      app.Qps,
			RealQpd:  appCounter.TotalCount,
			RealQps:  appCounter.QPS,
		})
	}

	return &types.APPListResponse{List: list, Total: 0}, nil
}
