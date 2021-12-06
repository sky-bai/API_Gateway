package global

import (
	"API_Gateway/api/internal/reverse_proxy/load_balance"
	"fmt"
	"strconv"
	"sync"
)

// 每个服务对应自己的负载均衡配置

var (
	LoadBalanceHandler    *LoadBalancer
	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1
)

type LoadBalancer struct {
	LoadBalanceMap   map[string]*LoadBalanceItem // 后端服务较多
	LoadBalanceSlice []*LoadBalanceItem          // 后端服务较少 减少锁的开销
	Locker           sync.RWMutex
}

type LoadBalanceItem struct {
	LoadBalance load_balance.LoadBalance
	ServiceName string
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		LoadBalanceMap:   map[string]*LoadBalanceItem{},
		LoadBalanceSlice: []*LoadBalanceItem{},
		Locker:           sync.RWMutex{},
	}
}

func init() {
	LoadBalanceHandler = NewLoadBalancer()
}

// GetLoadBalancer  这里会根据当前请求在数据库配置 去获取对应的负载均衡配置
func (lbr *LoadBalancer) GetLoadBalancer(service ServiceDetail) (load_balance.LoadBalance, error) {
	for _, item := range lbr.LoadBalanceSlice {
		if item.ServiceName == service.Info.ServiceName {
			return item.LoadBalance, nil
		}
	}
	schema := "http"
	if service.HTTPRule.NeedHttps == 1 {
		schema = "https"
	}

	prefix := ""
	if service.HTTPRule.RuleType == int64(HTTPRuleTypePrefixURL) {
		prefix = strconv.FormatInt(service.HTTPRule.RuleType, 10)
	}

	ipConf := make(map[string]string, 0)
	for key, ip := range service.LoadBalance.IpList {
		ipConf[string(ip)] = string(service.LoadBalance.WeightList[key])
	}
	mConf, err := load_balance.NewLoadBalanceCheckConf(fmt.Sprintf("%s://%s%s", schema, prefix), ipConf)
	if err != nil {
		return nil, err
	}
	// 把自己需要的负载均衡算法和负载均衡配置传入
	rb := load_balance.LoadBalanceFactorWithConf(load_balance.LbWeightRoundRobin, mConf)
	fmt.Println(rb)
	return nil, nil
}
