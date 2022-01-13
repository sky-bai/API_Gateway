package grpc_proxy_middleware

import (
	"API_Gateway/api/internal/global"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"github.com/dgrijalva/jwt-go"

	"strings"
)

// GrpcJwtAuthTokenMiddleware jwt auth token
func GrpcJwtAuthTokenMiddleware(serviceDetail *global.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error{
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error{
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}
		authToken:=""
		auths:=md.Get("authorization")
		if len(auths)>0{
			authToken = auths[0]
		}
		token:=strings.ReplaceAll(authToken,"Bearer ","")
		appMatched:=false
		if token!=""{
			claims,err:=JwtDecode(token)
			if err!=nil{
				return errors.WithMessage(err,"JwtDecode")
			}
			appList:=global.AppInfo.AppSlice
			for _,appInfo:=range appList{
				if appInfo.AppId==claims.Issuer{
					md.Set("app",o2json(appInfo))
					appMatched = true
					break
				}
			}
		}
		if serviceDetail.AccessControl.OpenAuth==1 && !appMatched{
			return errors.New("not match valid app")
		}
		if err := handler(srv, ss);err != nil {
			log.Printf("GrpcJwtAuthTokenMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}


func JwtDecode(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.JwtSignKey), nil
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
func o2json(o interface{}) string {
	marshal, err := jsoniter.Marshal(o)
	if err != nil {
		return "转换失败"
	}
	return string(marshal)
}
