package all

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AppTestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) AppTestLogic {
	return AppTestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppTestLogic) AppTest(req types.AppStatus) error {
	// todo: add your logic here and delete this line

	return nil
}
