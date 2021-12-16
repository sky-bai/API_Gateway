package admin

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AdminInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) AdminInfoLogic {
	return AdminInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminInfoLogic) AdminInfo() (*types.AdminInfoReponse, error) {
	// todo: add your logic here and delete this line

	return &types.AdminInfoReponse{}, nil
}
