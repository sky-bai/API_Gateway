package app

import (
	"API_Gateway/api/internal/middleware"
	"API_Gateway/model/ga_gateway_app"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"gopkg.in/go-playground/validator.v9"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddAppLogic {
	return AddAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type ErrorString struct {
	errMessage string
}

func (e *ErrorString) Error() string {
	return e.errMessage
}

// AddApp 添加租户
func (l *AddAppLogic) AddApp(req types.AddAppRequest) (*types.AppResponse, error) {
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

	// 2.验证app_id是否被占用
	_, err = l.svcCtx.GatewayAppModel.FindOneByAppId(req.AppID)
	if err != nil && err != sqlc.ErrNotFound {
		fmt.Println(err)
		return nil, errors.New("添加租户失败")
	}
	if err != sqlc.ErrNotFound {
		return nil, errors.New("该appId已被占用")
	}

	// 3.添加租户
	var app ga_gateway_app.GatewayApp
	app.AppId = req.AppID
	app.Name = req.Name
	app.Secret = req.Secret
	app.WhiteIps = req.WhiteIPS
	app.Qpd = req.Qpd
	app.Qps = req.Qps

	_, err = l.svcCtx.GatewayAppModel.Insert(app)
	if err != nil {
		return nil, errors.New("添加租户失败")
	}

	return &types.AppResponse{Message: "添加租户成功"}, nil

}
