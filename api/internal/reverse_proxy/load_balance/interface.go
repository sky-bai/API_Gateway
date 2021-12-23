package load_balance

type LoadBalance interface {
	Add(...string) error // 添加下游服务器节点
	Get(string) (string, error)
	Update()
}

type LoadBalanceConfInterface interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}
