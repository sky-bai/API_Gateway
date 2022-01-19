package middleware

import (
	"API_Gateway/api/internal/global"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type HTTPJwtTokenMiddleware struct {
}

func NewHTTPJwtTokenMiddleware() *HTTPJwtTokenMiddleware {
	return &HTTPJwtTokenMiddleware{}
}

const (
	JwtSignKey = "my_sign_key"
	JwtExpires = 60 * 60
)

func (m *HTTPJwtTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 1.这里是前面匹配到的http服务
		service := r.Context().Value("serviceInfo")
		serviceInfo := service.(*global.ServiceDetail)

		// decode jwt token
		// app_id 与  app_list 取得 appInfo
		// appInfo 放到 gin.context

		token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
		//fmt.Println("token",token)
		appMatched := false
		if token != "" {
			claims, err := JwtDecode(token)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			//fmt.Println("claims.Issuer",claims.Issuer)

			appList := global.AppInfo.AppSlice
			//fmt.Println("appList", appList)
			for _, appInfo := range appList {
				if appInfo.AppId == claims.Issuer {
					ctx := context.WithValue(r.Context(), "app", appInfo)
					next(w, r.WithContext(ctx))
					appMatched = true
					return
				}
			}
		}
		fmt.Println("--", serviceInfo.AccessControl.OpenAuth)
		fmt.Println(appMatched)
		if serviceInfo.AccessControl.OpenAuth == 1 && !appMatched {
			w.Write([]byte("该请求未授权"))
			return
		}

		next(w, r)
	}
}

func JwtDecode(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSignKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("token is not jwt.StandardClaims")
	}
}
