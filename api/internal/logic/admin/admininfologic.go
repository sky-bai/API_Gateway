package admin

import (
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"context"
	"encoding/json"

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

// AdminInfo 获取用户信息
func (l *AdminInfoLogic) AdminInfo() (interface{}, error) {

	value := l.ctx.Value("userId")
	userId, _ := value.(json.Number).Int64()

	adminInfo, err := l.svcCtx.GatewayAdminModel.FindOne(userId)
	if err != nil {
		return nil, err
	}

	out := types.AdminInfoReponse{
		ID:           int(adminInfo.Id),
		Name:         adminInfo.UserName,
		LoginTime:    int(adminInfo.CreateAt.Unix()),
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4670-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}

	return out, nil
}
