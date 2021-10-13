package admin

import (
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"API_Gateway/util"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/tal-tech/go-zero/core/logx"
	"time"
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

// 用户登陆
func (l *AdminLoginLogic) AdminLogin(req types.LoginRequest) (*types.LoginReponse, error) {
	validate := validator.New()
	// 通过validate去验证结构体
	err := validate.Struct(req)
	if err != nil {
		return nil, err
	}

	// 从数据库查找密码
	adminInfo, err := l.svcCtx.GatewayAdminModel.FindOneByUserName(req.UserName)
	if err != nil {
		fmt.Println("查询数据库err :", err)
		return nil, util.NewErrCode(util.UserNotFound)
	}

	// 匹配密码
	if req.Password == adminInfo.Password {
		return nil, util.NewErrCode(util.UserErrpPwd)
	}

	// ---start---
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, adminInfo.Id)
	if err != nil {
		return nil, util.NewErrCode(util.UserTokenFailSet)
	}
	// 生成了用户token

	return &types.LoginReponse{
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (l *AdminLoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
