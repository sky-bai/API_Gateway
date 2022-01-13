package app

import (
	"API_Gateway/api/internal/manager"
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

const (
	JwtSignKey = "my_sign_key"
	JwtExpires = 60 * 60
)

type GetTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetTokenLogic {
	return GetTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenLogic) GetToken(req types.GetTokenRequest) (*types.GetTokenResponse, error) {

	//  取出 app_id secret
	//  生成 app_list
	//  匹配 app_id
	//  基于 jwt生成token
	//  生成 output

	parts := l.ctx.Value("info").([]string)
	fmt.Println(parts)

	//fmt.Println("kkkkkkkkkkk", manager.AppHandler.GetAppList())
	manager.AppHandler.LoadOnce()

	appList := manager.AppHandler.AppSlice
	//fmt.Println("2222222223", appList)

	for _, appInfo := range appList {
		//fmt.Println("111", appInfo.AppId)
		//fmt.Println("111", appInfo.Secret)
		if appInfo.AppId == parts[0] && appInfo.Secret == parts[1] {
			claims := jwt.StandardClaims{
				Issuer:    appInfo.AppId,
				ExpiresAt: time.Now().Add(JwtExpires * time.Second).In(time.Local).Unix(),
			}
			token, err := JwtEncode(claims)
			if err != nil {
				return nil, errors.New("用户名或密码格式错误111")
			}
			return &types.GetTokenResponse{
				ExpiresIn:   JwtExpires,
				TokenType:   "Bearer",
				AccessToken: token,
				Scope:       "read_write"}, nil
		}
	}
	return nil, errors.New("没有该用户")
}

func JwtEncode(claims jwt.StandardClaims) (string, error) {
	mySigningKey := []byte(JwtSignKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}
