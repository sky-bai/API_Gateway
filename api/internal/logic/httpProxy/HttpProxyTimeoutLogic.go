package httpProxy

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type HttpProxyTimeoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpProxyTimeoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) HttpProxyTimeoutLogic {
	return HttpProxyTimeoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpProxyTimeoutLogic) HttpProxyTimeout() (resp *types.PingReponse, err error) {
	// todo: add your logic here and delete this line

	return
}
