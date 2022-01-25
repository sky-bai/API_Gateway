package global

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var TransportHandler *Transport

func init() {
	TransportHandler = NewTransport()
}

// TransportItem 每一个服务对应一个连接池
type TransportItem struct {
	Transport   *http.Transport
	ServiceName string
}
type Transport struct {
	TransportMap   map[string]*TransportItem
	TransportSlice []*TransportItem
	Locker         sync.RWMutex
}

func NewTransport() *Transport {
	return &Transport{
		TransportMap:   map[string]*TransportItem{},
		TransportSlice: []*TransportItem{},
		Locker:         sync.RWMutex{},
	}
}

// 负载均衡里面有连接池超时的设置

// GetTrans 为这次请求配置对应的连接池
func (t *Transport) GetTrans(service *ServiceDetail) (*http.Transport, error) {

	// 1.判断这次请求的服务是否已经有连接池
	//for _, transportItem := range t.TransportSlice {
	//	if transportItem.ServiceName == service.Info.ServiceName {
	//		return transportItem.Transport, nil
	//	}
	//}
	//fmt.Println("service.LoadBalance.UpstreamHeaderTimeout",service.LoadBalance.UpstreamHeaderTimeout)
	// 2.如果没有则创建一个连接池
	trans := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(service.LoadBalance.UpstreamConnectTimeout) * time.Second, // 连接超时
		}).DialContext,
		MaxIdleConns:          int(service.LoadBalance.UpstreamMaxIdle),                               // 最大空闲连接数  0 表示没有限制 保持 keep-alive 的数量
		IdleConnTimeout:       time.Duration(service.LoadBalance.UpstreamIdleTimeout) * time.Second,   // 空闲连接超时  这个连接在关闭之前保持空闲的最长时间
		ResponseHeaderTimeout: time.Duration(service.LoadBalance.UpstreamHeaderTimeout) * time.Second} // 读取返回头的超时

	// 3.将已有服务和对应的配置放入map
	TranItem := &TransportItem{
		Transport:   trans,
		ServiceName: service.Info.ServiceName}
	t.TransportSlice = append(t.TransportSlice, TranItem)

	//fmt.Println("设置的获取header头超时时间",trans.ResponseHeaderTimeout)

	t.Locker.Lock()
	defer t.Locker.Unlock()
	t.TransportMap[service.Info.ServiceName] = TranItem

	return trans, nil
}
