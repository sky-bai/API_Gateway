package httpProxy

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type HttpProxyWhiteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpProxyWhiteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) HttpProxyWhiteListLogic {
	return HttpProxyWhiteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpProxyWhiteListLogic) HttpProxyWhiteList() (resp *types.PingReponse, err error) {
	// todo: add your logic here and delete this line

	return
}
