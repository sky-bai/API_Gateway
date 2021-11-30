package admin

import (
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"context"

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

// 获取用户信息
func (l *AdminInfoLogic) AdminInfo() (*types.AdminInfoReponse, error) {

	return &types.AdminInfoReponse{}, nil
}
