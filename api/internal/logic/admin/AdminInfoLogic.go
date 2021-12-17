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

// AdminInfo 获取管理员信息
func (l *AdminInfoLogic) AdminInfo() (*types.AdminInfoReponse, error) {
	userId := 1
	adminInfo, err := l.svcCtx.GatewayAdminModel.FindOne(int64(userId))
	if err != nil {
		return nil, err
	}

	return &types.AdminInfoReponse{
		ID:           int(adminInfo.Id),
		Name:         adminInfo.UserName,
		LoginTime:    int(adminInfo.CreateAt.Unix()),
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4670-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}, nil

}
