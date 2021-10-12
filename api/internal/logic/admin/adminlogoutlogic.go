package admin

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

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

	return &types.LogOutReponse{}, nil
}
