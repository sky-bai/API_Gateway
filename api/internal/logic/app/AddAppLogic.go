package app

import (
	"context"

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

func (l *AddAppLogic) AddApp(req types.AddAppRequest) (*types.AppResponse, error) {
	// todo: add your logic here and delete this line

	return &types.AppResponse{}, nil
}
