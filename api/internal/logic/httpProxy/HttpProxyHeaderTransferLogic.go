package httpProxy

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type HttpProxyHeaderTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpProxyHeaderTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) HttpProxyHeaderTransferLogic {
	return HttpProxyHeaderTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpProxyHeaderTransferLogic) HttpProxyHeaderTransfer() (resp *types.PingReponse, err error) {
	// todo: add your logic here and delete this line

	return
}
