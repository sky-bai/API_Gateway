package global

import (
	"API_Gateway/reverse_proxy/load_balance"
	"fmt"
	"sync"
)

// 每个服务对应自己的负载均衡配置

var (
	LoadBalanceHandler    *LoadBalancer
	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1
)

type LoadBalancer struct {
	LoadBalanceMap   map[string]load_balance.LoadBalance // 后端服务较多
	LoadBalanceSlice []load_balance.LoadBalance          // 后端服务较少 减少锁的开销
	Locker           sync.RWMutex
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		LoadBalanceMap:   make(map[string]load_balance.LoadBalance),
		LoadBalanceSlice: make([]load_balance.LoadBalance, 0),
		Locker:           sync.RWMutex{},
	}
}

func init() {
	LoadBalanceHandler = NewLoadBalancer()
}
func (lbr *LoadBalancer) GetLoadBalancer(service ServiceDetail) (load_balance.LoadBalance, error) {
	schema := "http"
	if service.HTTPRule.NeedHttps == 1 {
		schema = "https"
	}

	prefix := ""
	if service.HTTPRule.RuleType == int64(HTTPRuleTypePrefixURL) {
		prefix = string(service.HTTPRule.RuleType)
	}

	mConf, err := load_balance.NewLoadBalanceCheckConf(fmt.Sprintf("%s://%s%s", schema, prefix), map[string]string{"127.0.0.1:2003": "20", "127.0.0.1:2004": "20"})
	if err != nil {
		panic(err)
	}
	rb := load_balance.LoadBalanceFactoryWithConf(load_balance.LbWeightRoundRobin, mConf)
	fmt.Println(rb)
	return nil, nil
}
