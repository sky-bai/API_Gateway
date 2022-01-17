package middleware

import (
	"API_Gateway/api/internal/global"
	"API_Gateway/model/ga_gateway_app"
	"fmt"
	"net/http"
)

type HTTPJwtFlowCountMiddleware struct {
}

func NewHTTPJwtFlowCountMiddleware() *HTTPJwtFlowCountMiddleware {
	return &HTTPJwtFlowCountMiddleware{}
}

func (m *HTTPJwtFlowCountMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		appList := r.Context().Value("app")
		fmt.Println("121312", appList)
		// 获取不到可能是没有开始jwt权限 就执行下游中间件
		if appList == nil {
			next(w, r)
			return
		}
		appInfo := appList.(*ga_gateway_app.GatewayApp)

		appCounter, err := global.FlowCounterHandler.GetCounter(FlowAppPrefix + appInfo.AppId)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		appCounter.Increase()
		// 如果超过了租户设置的请求数 就丢弃请求
		fmt.Println("appInfo.Qpd", appInfo.Qpd)
		fmt.Println("appCounter.TotalCount", appCounter.TotalCount)

		if appInfo.Qpd > 0 && appCounter.TotalCount > appInfo.Qpd {
			w.Write([]byte(fmt.Sprintf("租户日请求量限流 limit:%v current:%v", appInfo.Qpd, appCounter.TotalCount)))
			return
		}
		next(w, r)
	}
}
