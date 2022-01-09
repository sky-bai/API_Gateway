package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

type AppTokenMiddleware struct {
}

func NewAppTokenMiddleware() *AppTokenMiddleware {
	return &AppTokenMiddleware{}
}

func (m *AppTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		splits := strings.Split(r.Header.Get("Authorization"), " ")
		if len(splits) != 2 {
			w.Write([]byte(fmt.Sprintf("用户名或密码格式错误")))
			return
		}
		// 2. token解码
		appSecret, err := base64.StdEncoding.DecodeString(splits[1])
		if err != nil {
			w.Write([]byte(fmt.Sprintf("解码appSecret错误")))
			return
		}
		fmt.Println("appSecret:", string(appSecret))
		parts := strings.Split(string(appSecret), ":")
		if len(parts) != 2 {
			w.Write([]byte(fmt.Sprintf("用户名或密码格式错误")))
			return
		}
		ctx := context.WithValue(r.Context(), "info", parts)
		next(w, r.WithContext(ctx))
	}
}
