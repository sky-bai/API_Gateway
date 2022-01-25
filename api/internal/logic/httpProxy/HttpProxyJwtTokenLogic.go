package httpProxy

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type HttpProxyJwtTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpProxyJwtTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) HttpProxyJwtTokenLogic {
	return HttpProxyJwtTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpProxyJwtTokenLogic) HttpProxyJwtToken() (resp *types.PingReponse, err error) {
	// todo: add your logic here and delete this line

	return
}
