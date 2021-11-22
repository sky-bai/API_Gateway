package load_balance

type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)
	Update()
}

// Observer 客户端需要一个监听者  监听服务列表变化然后更新服务列表
type Observer interface {
	Update()
	// 提供一个有负载均衡配置的观察者更新接口  下面是使用实例
}

// LoadBalanceZkConfInterface 配置主题   负载均衡配置绑定一个观察者 这里是服务端注册发现 需要有客户端的观察者
type LoadBalanceZkConfInterface interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

// 也就是说我将自身改变服务器列表的时候 我能够自身update服务器的方法 不管你是什么
// LoadBalanceConf 是可以监听服务器列表的实体和更新配置的实体

// 也就是说我将两个实体抽象出来 确定了他们的方法
