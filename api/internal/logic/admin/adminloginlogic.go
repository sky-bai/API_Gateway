package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strconv"

	"api/internal/svc"
	"api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AdminLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) AdminLoginLogic {
	return AdminLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLoginLogic) AdminLogin(req types.LoginRequest) (*types.LoginReponse, error) {
	validate := validator.New()
	// 通过validate去验证结构体
	err := validate.Struct(req)
	if err != nil {
		return nil, err
	}
	// 1.获取到用户id
	logx.Infof("userId: %v", l.ctx.Value("userId")) // 这里的key和生成jwt token时传入的key一致
	value := l.ctx.Value("userId")
	i, _ := value.(json.Number).Int64()
	strInt64 := strconv.FormatInt(i, 10)
	userID, _ := strconv.Atoi(strInt64)
	fmt.Println(userID)

	return &types.LoginReponse{}, nil
}
