package middleware

import (
	"fmt"
	"net/http"
)

type HTTPAccessModeMiddleware struct {
}

func NewHTTPAccessModeMiddleware() *HTTPAccessModeMiddleware {
	return &HTTPAccessModeMiddleware{}
}

// 这个中间件去判断接入过来的请求是否后端有这个服务
func (m *HTTPAccessModeMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("http proxy ", r)

		next(w, r)
	}
}
