package httpsProxy

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type HttpsProxyPingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpsProxyPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) HttpsProxyPingLogic {
	return HttpsProxyPingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpsProxyPingLogic) HttpsProxyPing() (*types.PingReponse, error) {

	return &types.PingReponse{Message: "https_proxy 可以跑通"}, nil
}
