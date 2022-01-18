package middleware

import (
	"API_Gateway/api/internal/global"
	"fmt"
	"net/http"
	"strings"
)

type HTTPWhiteListMiddleware struct {
}

func NewHTTPWhiteListMiddleware() *HTTPWhiteListMiddleware {
	return &HTTPWhiteListMiddleware{}
}

// Handle 设置IP白名单
func (m *HTTPWhiteListMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1.这里是前面匹配到的http服务
		service := r.Context().Value("serviceInfo")
		serviceInfo := service.(*global.ServiceDetail)
		fmt.Println("r.RemoteAddr", r.RemoteAddr)
		forwarded := r.Header.Get("X-Forwarded-For")
		fmt.Println("X-Forwarded-For", forwarded)

		ip := r.Header.Get("X-Real-IP")
		fmt.Println("获取到的ip:", ip)
		if ip == "" {
			ip = r.Header.Get("X-Forwarded-For")
		}

		var iplist []string
		if serviceInfo.AccessControl.WhiteList != "" {
			iplist = strings.Split(serviceInfo.AccessControl.WhiteList, ",")
		}
		fmt.Println("白名单iplist:", iplist)
		if serviceInfo.AccessControl.OpenAuth == 1 && len(iplist) > 0 {
			if !InStringSlice(iplist, r.Header.Get("X-Real-IP")) {
				w.Write([]byte("当前请求IP不在白名单中"))
				return
			}
		}

		next(w, r)
	}
}

func InStringSlice(slice []string, str string) bool {
	fmt.Println("slice:", slice)
	fmt.Println("str:", str)
	for _, item := range slice {
		if str == item {
			return true
		}
	}
	return false
}
