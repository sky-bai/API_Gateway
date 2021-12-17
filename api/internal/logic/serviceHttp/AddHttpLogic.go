package serviceHttp

import (
	"context"

	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddHttpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddHttpLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddHttpLogic {
	return AddHttpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddHttpLogic) AddHttp(req types.AddHTTPResquest) (*types.HttpReponse, error) {
	// todo: add your logic here and delete this line

	return &types.HttpReponse{}, nil
}
