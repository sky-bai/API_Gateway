package middleware

import (
	"API_Gateway/api/internal/global"
	"API_Gateway/model/ga_gateway_app"
	"fmt"
	"net/http"
)

type HTTP_Jwt_Flow_LimitMiddleware struct {
}

func NewHTTP_Jwt_Flow_LimitMiddleware() *HTTP_Jwt_Flow_LimitMiddleware {
	return &HTTP_Jwt_Flow_LimitMiddleware{}
}

func (m *HTTP_Jwt_Flow_LimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 1.获取租户信息
		appList := r.Context().Value("app")
		// 获取不到可能是没有开始jwt权限 就执行下游中间件
		if appList == nil {
			next(w, r)
			return
		}
		appInfo := appList.(*ga_gateway_app.GatewayApp)

		if appInfo.Qps > 0 {
			clientLimiter, err := global.FlowLimiterHandler.GetLimiter(
				global.FlowAppPrefix+appInfo.AppId+"_"+r.RemoteAddr,
				float64(appInfo.Qps))
			if err != nil {
				w.Write([]byte(fmt.Sprintf("%s", err)))
				return
			}
			if !clientLimiter.Allow() {
				w.Write([]byte(fmt.Sprintf("%v flow limit %v", r.RemoteAddr, appInfo.Qps)))
				return
			}
		}
		next(w, r)
	}
}
