package app

import (
	"context"
	"errors"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AppDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) AppDetailLogic {
	return AppDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AppDetail 获取租户详细信息
func (l *AppDetailLogic) AppDetail(req types.AppDetailRequest) (interface{}, error) {

	serviceList, err := l.svcCtx.GatewayAppModel.FindOne(req.ID)
	if err != nil {
		return nil, errors.New("获取租户详细信息失败")
	}

	return serviceList, nil
}
