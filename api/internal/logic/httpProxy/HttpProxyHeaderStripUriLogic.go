package httpProxy

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type HttpProxyHeaderStripUriLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpProxyHeaderStripUriLogic(ctx context.Context, svcCtx *svc.ServiceContext) HttpProxyHeaderStripUriLogic {
	return HttpProxyHeaderStripUriLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpProxyHeaderStripUriLogic) HttpProxyHeaderStripUri() (resp *types.PingReponse, err error) {
	// todo: add your logic here and delete this line

	return
}
