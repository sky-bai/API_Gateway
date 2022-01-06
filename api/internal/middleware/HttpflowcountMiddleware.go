package middleware

import (
	"API_Gateway/api/internal/global"
	"fmt"
	"net/http"
)

type HTTPFlowCountMiddleware struct {
}

func NewHTTPFlowCountMiddleware() *HTTPFlowCountMiddleware {
	return &HTTPFlowCountMiddleware{}
}

const (
	FlowTotal         = "flow_total"
	FlowServicePrefix = "flow_service_"
	FlowAppPrefix     = "flow_app_"
)

func (m *HTTPFlowCountMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 1.这里是前面匹配到的http服务
		service := r.Context().Value("serviceInfo")
		serviceInfo := service.(*global.ServiceDetail)

		// 2.获取本次请求对应流量控制器

		// 1.全站
		totalCounter, err := global.FlowCounterHandler.GetCounter(FlowTotal)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("get total flow counter error: %v", err)))
			return
		}
		totalCounter.Increase()

		// 2.服务
		serviceCounter, err := global.FlowCounterHandler.GetCounter(FlowServicePrefix + serviceInfo.Info.ServiceName)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("get service flow counter error: %v", err)))
			return
		}
		serviceCounter.Increase()

		// 3.租户
		appCounter, err := global.FlowCounterHandler.GetCounter(FlowAppPrefix)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("get total flow counter error: %v", err)))
			return
		}
		appCounter.Increase()

		next(w, r)
	}
}
