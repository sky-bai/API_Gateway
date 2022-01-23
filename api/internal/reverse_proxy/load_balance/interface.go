package load_balance

type LoadBalance interface {
	Add(...string) error // 添加下游服务器节点
	// Get 2.获取下游服务器
	Get(string) (string, error)
	// Update 更新服务器节点
	Update()
}

// LoadBalanceConfInterface 维护下游服务器节点列表的状态
type LoadBalanceConfInterface interface {
	// Attach 1.添加负载均衡观察者
	Attach(o Observer)
	// GetConf 2.获取下游服务器列表
	GetConf() []string
	// WatchConf 3.主动监听下游服务器列表变化
	WatchConf()
	// UpdateConf 4.通知观察者更新服务器列表
	UpdateConf(conf []string)
}
