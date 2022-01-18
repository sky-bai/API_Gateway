package middleware

import (
	"API_Gateway/api/internal/global"
	"API_Gateway/api/internal/reverse_proxy"
	"fmt"
	"net/http"
)

type HTTPReverseProxyMiddleware struct {
}

func NewHTTPReverseProxyMiddleware() *HTTPReverseProxyMiddleware {
	return &HTTPReverseProxyMiddleware{}
}

// Handle 对接后端接口 基于请求信息 配置代理 反向代理的中间件
func (m *HTTPReverseProxyMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1.这里是前面匹配到的http服务
		service := r.Context().Value("serviceInfo")
		serviceInfo := service.(*global.ServiceDetail)
		fmt.Println("reverse proxy")

		//fmt.Println("请求对应的下游服务器:", serviceInfo.LoadBalance.IpList)
		//fmt.Println("请求对应的下游服务器权重:", serviceInfo.LoadBalance.WeightList)

		// 2.根据服务在数据库配置的下游服务器列表ipList和负载类型loadType 获得一个负载均衡器
		loadBalancer, err := global.LoadBalanceHandler.GetLoadBalancer(*serviceInfo) // 负载均衡只维护下游服务器地址 不维护path
		//fmt.Println("loadBalancer:", Obj2Json(loadBalancer))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// 3.每个服务有自己对应的连接池 去设置每个服务的连接池 反向代理有这个连接池这个设置
		transport, err := global.TransportHandler.GetTrans(serviceInfo)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// 4.设置代理
		reverse_proxy.NewLoadBalanceReverseProxy(r, loadBalancer, transport).ServeHTTP(w, r)
	}
}
