package global

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var TransportHandler *http.Transport

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

func (t *Transport) GetTrans(service *ServiceDetail) (*http.Transport, error) {
	for _, transportItem := range t.TransportSlice {
		if transportItem.ServiceName == service.Info.ServiceName {
			return transportItem.Transport, nil
		}
	}
	trans := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(service.LoadBalance.UpstreamConnectTimeout),
		}).DialContext,
		MaxIdleConns:          int(service.LoadBalance.UpstreamMaxIdle),
		IdleConnTimeout:       time.Duration(service.LoadBalance.UpstreamIdleTimeout),
		ResponseHeaderTimeout: time.Duration(service.LoadBalance.UpstreamHeaderTimeout)}

	TranItem := &TransportItem{
		Transport:   trans,
		ServiceName: service.Info.ServiceName}
	t.TransportSlice = append(t.TransportSlice, TranItem)

	t.Locker.Lock()
	defer t.Locker.Unlock()
	t.TransportMap[service.Info.ServiceName] = TranItem

	return trans, nil
}
