package httpProxy

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type HttpProxyUrlRewriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpProxyUrlRewriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) HttpProxyUrlRewriteLogic {
	return HttpProxyUrlRewriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpProxyUrlRewriteLogic) HttpProxyUrlRewrite() (resp *types.PingReponse, err error) {
	// todo: add your logic here and delete this line

	return
}
