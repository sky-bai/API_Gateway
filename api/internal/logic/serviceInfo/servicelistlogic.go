package serviceInfo

import (
	"context"
	"fmt"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ServiceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServiceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) ServiceListLogic {
	return ServiceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServiceListLogic) ServiceList(req types.ServiceListResquest) (interface{}, error) {
	// todo: add your logic here and delete this line
	fmt.Println("sdfsdfsd")
	dataLike, err := l.svcCtx.GatewayServiceInfoModel.FindDataLike(req.Info, req.PageSize, req.PageNo)
	if err != nil {
		return nil, err
	}
	return &dataLike, nil
}
