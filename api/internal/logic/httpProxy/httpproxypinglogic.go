package httpProxy

import (
	"API_Gateway/api/internal/global"
	"context"
	"errors"
	"fmt"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

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
	fmt.Println("👌👌👌", l.ctx.Value("serviceInfo"))
	detail, exist := l.ctx.Value("serviceInfo").(global.ServiceDetail)
	if !exist {
		return nil, errors.New("未匹配到服务信息")
	}
	fmt.Println("服务信息", detail.HTTPRule)
	return &types.PingReponse{
		Message: "http_proxy 可以跑通",
	}, nil
}
