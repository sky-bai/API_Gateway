package middleware

import (
	"API_Gateway/api/internal/global"
	"fmt"
	"net/http"
)

type HTTPFlowLimitMiddleware struct {
}

func NewHTTPFlowLimitMiddleware() *HTTPFlowLimitMiddleware {
	return &HTTPFlowLimitMiddleware{}
}

func (m *HTTPFlowLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1.这里是前面匹配到的http服务
		service := r.Context().Value("serviceInfo")
		serviceInfo := service.(*global.ServiceDetail)

		// 等于 0 表示不限流 服务端限流
		if serviceInfo.AccessControl.ServiceFlowLimit != 0 {
			// 根据 serviceInfo.AccessControl.ServiceFlowLimit 限流数
			// 1.构建限流器
			serviceLimiter, err := global.FlowLimiterHandler.GetLimiter(
				global.FlowServicePrefix+serviceInfo.Info.ServiceName,
				float64(serviceInfo.AccessControl.ServiceFlowLimit))
			if err != nil {
				w.Write([]byte(fmt.Sprintf("%s", err)))
				return
			}
			// 2.如果没有令牌就返回
			if !serviceLimiter.Allow() {
				w.Write([]byte(fmt.Sprintf("\"当前请求已达到峰值，服务端开启限流 service flow limit %v", serviceInfo.AccessControl.ServiceFlowLimit)))
				return
			}
		}

		// 对客户端进行限流 获取到请求到ip
		if serviceInfo.AccessControl.ClientipFlowLimit > 0 {
			//fmt.Println("=====",r.RemoteAddr)
			clientLimiter, err := global.FlowLimiterHandler.GetLimiter(
				global.FlowServicePrefix+serviceInfo.Info.ServiceName+"_"+r.RemoteAddr,
				float64(serviceInfo.AccessControl.ClientipFlowLimit))
			if err != nil {
				w.Write([]byte(fmt.Sprintf("%s", err)))
				return
			}
			fmt.Println("")
			if !clientLimiter.Allow() {
				w.Write([]byte(fmt.Sprintf("请求的IP %v flow limit %v", r.RemoteAddr, serviceInfo.AccessControl.ClientipFlowLimit)))
				return
			}
		}

		next(w, r)
	}
}
