package load_balance

type LbType int

const (
	LbRandom LbType = iota
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
)

// LoadBalanceFactoryWithConf 我把你想要的实例都创建出来
func LoadBalanceFactoryWithConf(lbType LbType, lbConf LoadBalanceZkConfInterface) LoadBalance {
	switch lbType {
	case LbRandom:
		lb := &RandomBalance{}
		// 1.为负载均衡配置zk配置
		lb.SetConf(lbConf)
		// 2.为zk配置它的观察者
		lbConf.Attach(lb)
		return lb
	case LbRoundRobin:
		lb := &RoundRobinBalance{}
		lb.SetConf(lbConf)
		lbConf.Attach(lb)
		return lb
	case LbWeightRoundRobin:
		lb := &WeightRoundRobinBalance{}
		lb.SetConf(lbConf)
		lbConf.Attach(lb)
		return lb
	case LbConsistentHash:
		lb := &ConsistentHashBalance{}
		lb.SetConf(lbConf)
		lbConf.Attach(lb)
		return lb
	default:
		lb := &RandomBalance{}
		lb.SetConf(lbConf)
		lbConf.Attach(lb)
		return lb
	}

}
