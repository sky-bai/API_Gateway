package middleware

import (
	"API_Gateway/api/internal/global"
	"net/http"
	"strings"
)

type HeaderTransferMiddleware struct {
}

func NewHeaderTransferMiddleware() *HeaderTransferMiddleware {
	return &HeaderTransferMiddleware{}
}

// Handle 修改请求重写层的头部信息
func (m *HeaderTransferMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		service := r.Context().Value("serviceInfo")
		serviceInfo := service.(*global.ServiceDetail)

		for _, value := range strings.Split(serviceInfo.HTTPRule.HeaderTransfor, ",") {
			item := strings.Split(value, " ")
			if len(item) != 3 {
				continue
			}
			if item[0] == "add" || item[0] == "edit" {
				r.Header.Set(item[1], item[2])
			}
			if item[0] == "del" {
				r.Header.Del(item[1])
			}

		}
		next(w, r)
	}
}
