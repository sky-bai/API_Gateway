package httpProxy

import (
	"context"

	"API_Gateway/api/interna/svc"
	"API_Gateway/api/interna/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type HttpProxyPingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpProxyPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) HttpProxyPingLogic {
	return HttpProxyPingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpProxyPingLogic) HttpProxyPing() (*types.PingReponse, error) {

	return &types.PingReponse{
		Message: "http_proxy 可以跑通",
	}, nil
}
