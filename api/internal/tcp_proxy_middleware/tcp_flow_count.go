package tcp_proxy_middleware

import (
	"API_Gateway/api/internal/global"
)

func TCPFlowCountMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(global.ServiceDetail)

		// 统计项 1 全站 2 服务 3 租户
		totalCounter, err := global.FlowCounterHandler.GetCounter(global.FlowTotal)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		totalCounter.Increase()

		// 服务
		serviceCounter, err := global.FlowCounterHandler.GetCounter(global.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		serviceCounter.Increase()
		c.Next()
	}
}
