package app

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

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

func (l *AppDetailLogic) AppDetail(req types.AppDetailRequest) (*types.AppDetailResponse, error) {
	fmt.Println("7889512223", req.ID)
	appDetail, err := l.svcCtx.GatewayAppModel.FindOne(int64(req.ID))
	if err != nil {
		return nil, errors.New("获取租户详细信息失败")
	}

	return &types.AppDetailResponse{ID: appDetail.Id,
		AppID:     appDetail.AppId,
		WhiteIPS:  appDetail.WhiteIps,
		Name:      appDetail.Name,
		Secret:    appDetail.Secret,
		Qpd:       appDetail.Qpd,
		Qps:       appDetail.Qps,
		CreatedAt: appDetail.CreateTime.Format("2006-06-06"),
		UpdatedAt: appDetail.UpdateTime.Format("2006-06-06")}, nil
}
