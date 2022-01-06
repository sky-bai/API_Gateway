package httpProxy

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type HttpProxyFlowCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpProxyFlowCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) HttpProxyFlowCountLogic {
	return HttpProxyFlowCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpProxyFlowCountLogic) HttpProxyFlowCount() (resp *types.PingReponse, err error) {
	// todo: add your logic here and delete this line

	return
}
