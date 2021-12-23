package app

import (
	"API_Gateway/api/internal/middleware"
	"context"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpdateAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateAppLogic {
	return UpdateAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateApp 更新租户
func (l *UpdateAppLogic) UpdateApp(req types.UpdateAppRequest) (*types.AppResponse, error) {
	// 1.校验参数
	errMessage := ErrorString{errMessage: ""}
	err := middleware.ValidatorHandler.Validate.Struct(&req)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, errValue := range errs.Translate(middleware.ValidatorHandler.Translate) {
			errMessage.errMessage += " " + errValue
		}
		return nil, &errMessage
	}

	// 2.验证是否存在该服务
	appInfo, err := l.svcCtx.GatewayAppModel.FindOne(req.ID)
	if err != nil {
		return nil, errors.New("查询租户信息失败")
	}
	if appInfo.Id < 0 {
		return nil, errors.New("租户不存在")
	}

	// 3.更新租户信息
	appInfo.Name = req.Name
	appInfo.AppId = req.AppID
	appInfo.WhiteIps = req.WhiteIPS
	appInfo.Qps = req.Qps
	appInfo.Qpd = req.Qpd
	appInfo.Secret = req.Secret
	err = l.svcCtx.GatewayAppModel.Update(*appInfo)
	if err != nil {
		return nil, errors.New("更新租户信息失败")
	}

	return &types.AppResponse{Message: "更新租户信息成功"}, nil

}
