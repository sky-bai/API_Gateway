package admin

import (
	"context"

	"API_Gateway/api/interna/svc"
	"API_Gateway/api/interna/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AdminLogOutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminLogOutLogic(ctx context.Context, svcCtx *svc.ServiceContext) AdminLogOutLogic {
	return AdminLogOutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 管理员退出
func (l *AdminLogOutLogic) AdminLogOut() (*types.LogOutReponse, error) {
	// 如果清除token

	return &types.LogOutReponse{}, nil
}
