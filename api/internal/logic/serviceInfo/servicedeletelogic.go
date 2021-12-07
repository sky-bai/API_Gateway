package serviceInfo

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ServiceDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServiceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) ServiceDeleteLogic {
	return ServiceDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ServiceDelete 服务删除
func (l *ServiceDeleteLogic) ServiceDelete(req types.ServiceResquest) (*types.CommonReponse, error) {
	err := l.svcCtx.GatewayServiceInfoModel.Delete(req.ID)
	if err != nil {
		return nil, err
	}
	return &types.CommonReponse{Msg: "删除成功"}, nil
}
