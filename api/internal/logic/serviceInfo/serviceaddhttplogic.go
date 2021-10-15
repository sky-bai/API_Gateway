package serviceInfo

import (
	"context"
	"fmt"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ServiceAddHttpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServiceAddHttpLogic(ctx context.Context, svcCtx *svc.ServiceContext) ServiceAddHttpLogic {
	return ServiceAddHttpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 添加http服务
func (l *ServiceAddHttpLogic) ServiceAddHttp(req types.AddHTTPResquest) (*types.CommonReponse, error) {

	fmt.Println(req)
	return &types.CommonReponse{}, nil
}
