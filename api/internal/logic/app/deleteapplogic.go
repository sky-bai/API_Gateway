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

type DeleteAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteAppLogic {
	return DeleteAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteApp 删除该租户
func (l *DeleteAppLogic) DeleteApp(req types.DeleteAppRequest) (*types.AppResponse, error) {
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

	// 2.查看是否有该租户
	_, err = l.svcCtx.GatewayAppModel.FindOne(int64(req.ID))
	if err != nil {
		return nil, errors.New("查询该租户信息失败")
	}

	// 3.删除租户
	err = l.svcCtx.GatewayAppModel.Delete(int64(req.ID))
	if err != nil {
		return nil, errors.New("删除该租户失败")
	}

	return &types.AppResponse{Message: "删除该租户信息成功"}, nil
}
