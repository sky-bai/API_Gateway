package global

import (
	"API_Gateway/api/internal/reverse_proxy/load_balance"
	"fmt"
	"strings"
	"sync"
)

// 每个服务对应自己的负载均衡配置

var (
	LoadBalanceHandler    *LoadBalancer
	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1
	LoadTypeTCP           = 1
	LoadTypeGRPC          = 2
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

// 负载均衡的一个管理器 每个服务对应自己的负载均衡配置
func init() {
	LoadBalanceHandler = NewLoadBalancer()
}

// 最终目的是获得一个负载均衡器 先要创建一个负载均衡配置 然后去获取负载均衡器

// GetLoadBalancer  这里会根据当前请求在数据库配置 去获取对应的负载均衡配置
func (lbr *LoadBalancer) GetLoadBalancer(service ServiceDetail) (load_balance.LoadBalance, error) {

	// 1. 根据serviceName 获取对应的负载均衡配置
	for _, item := range lbr.LoadBalanceSlice {
		if item.ServiceName == service.Info.ServiceName {
			return item.LoadBalance, nil
		}
	}
	// 2.创建这个服务的format
	schema := "http://"
	if service.HTTPRule.NeedHttps == 1 {
		schema = "https://"
	}

	//prefix := ""
	//if service.HTTPRule.RuleType == int64(HTTPRuleTypePrefixURL) {
	//	prefix = strconv.FormatInt(service.HTTPRule.RuleType, 10)
	//	fmt.Println("prefix", prefix)
	//}
	if int(service.Info.LoadType) == LoadTypeTCP || int(service.Info.LoadType) == LoadTypeGRPC {
		schema = ""
	}

	// 3.这里维护一个服务器地址与权重的k v
	ipConf := make(map[string]string)
	//fmt.Println("下游服务器列表",service.LoadBalance.IpList)
	ipArrayStr := strings.Split(service.LoadBalance.IpList, ",")
	weightArray := strings.Split(service.LoadBalance.WeightList, ",")

	for key, ip := range ipArrayStr {
		ipConf[ip] = weightArray[key]
	}
	//fmt.Println("下游服务器与权重信息",ipConf)
	// 负载均衡需要下游服务器地址 ，这个根据这个服务获得在数据库保存的下游服务器地址

	//fmt.Println("前缀",fmt.Sprintf("%s%s", schema, "%s"))

	// 4.创建一个管理服务器列表的对象 负载均衡配置  获取一个负载均衡配置 并且主动探测下游服务器
	mConf, err := load_balance.NewLoadBalanceCheckConf(fmt.Sprintf("%s%s", schema, "%s"), ipConf)
	if err != nil {
		return nil, err
	}

	// 5.把自己需要的负载均衡算法和负载均衡配置传入  获得一个负载均衡器
	rb := load_balance.LoadBalanceFactorWithConf(load_balance.LbType(service.LoadBalance.RoundType), mConf)
	//fmt.Println(rb)
	//fmt.Println("本次请求对应的负载均衡类型",service.LoadBalance.RoundType)

	// 6.把这个负载均衡器和serviceName 存入map
	lbItem := &LoadBalanceItem{
		LoadBalance: rb,
		ServiceName: service.Info.ServiceName,
	}
	lbr.LoadBalanceSlice = append(lbr.LoadBalanceSlice, lbItem)
	lbr.Locker.Lock()
	defer lbr.Locker.Unlock()
	lbr.LoadBalanceMap[service.Info.ServiceName] = lbItem

	return rb, nil
}
