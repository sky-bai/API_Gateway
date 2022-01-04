package middleware

import (
	"API_Gateway/api/internal/global"
	"fmt"
	"net/http"
	"strings"
)

type HTTPBlackListMiddleware struct {
}

func NewHTTPBlackListMiddleware() *HTTPBlackListMiddleware {
	return &HTTPBlackListMiddleware{}
}

func (m *HTTPBlackListMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 1.这里是前面匹配到的http服务
		service := r.Context().Value("serviceInfo")
		serviceInfo := service.(*global.ServiceDetail)

		var whileIpList []string
		if serviceInfo.AccessControl.WhiteList != "" {
			whileIpList = strings.Split(serviceInfo.AccessControl.WhiteList, ",")
		}

		var blackIpList []string
		if serviceInfo.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceInfo.AccessControl.BlackList, ",")
		}
		fmt.Println("黑名单", blackIpList)
		fmt.Println("real-ip", r.Header.Get("X-Real-IP"))
		if serviceInfo.AccessControl.OpenAuth == 1 && len(whileIpList) == 0 && len(blackIpList) > 0 {
			if InStringSlice(blackIpList, r.Header.Get("X-Real-IP")) {
				next(w, r)
			} else {
				_, err := w.Write([]byte("Access Denied"))
				if err != nil {
					fmt.Println("write error", err)
					return
				}
			}
		}
	}
}
