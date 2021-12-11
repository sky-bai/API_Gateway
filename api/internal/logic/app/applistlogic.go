package app

import (
	"context"
	"errors"

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

// AppList 获取租户列表
func (l *AppListLogic) AppList(req types.AppListRequest) (interface{}, error) {

	appList, err := l.svcCtx.GatewayAppModel.GetServiceList(req.Info, req.PageNo, req.PageSize)
	if err != nil {
		return nil, errors.New("获取租户列表失败")
	}

	return appList, nil
}
