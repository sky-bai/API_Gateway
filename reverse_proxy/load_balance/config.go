package load_balance

import (
	"fmt"
	"github.com/e421083458/gateway_demo/proxy/zookeeper"
)

// 现在要实例化负载均衡zk配置 这里面要管理服务器列表 监听服务器列表的变化 通过配置zk服务器信息动态获取到服务器列表 可以获取多个观察者
// 将观察者 负载均衡方法和负载均衡配置绑定
// 这里LoadBalanceZkConf 实现了上面LoadBalanceConf的接口

type LoadBalanceZkConf struct {
	observers    []Observer // 接口 一个接口 相当于 抽象出一个可以更新的观察者 在使用的时候我只需要关心该接口所实现的方法
	format       string
	zkHosts      []string
	zkPath       string
	activeList   []string          // 可用IP列表
	confIPWeight map[string]string // 下游服务器和权重
}

// NewLoadBalanceZkConf 实例化负载均衡zk配置 并实时监听服务器列表的变化 如果有变化就会通知负载均衡调用方更新它们的服务器列表
func NewLoadBalanceZkConf(format, path string, zkHosts []string, conf map[string]string) (*LoadBalanceZkConf, error, *zookeeper.ZkManager) {
	zkManager := zookeeper.NewZkManager(zkHosts)
	err := zkManager.GetConnect()
	if err != nil {
		fmt.Println("zkManager.GetConnect err:", err)
		return nil, err, nil
	}
	//defer zkManager.Close()

	zkList, err := zkManager.GetServerListByPath(path)
	fmt.Println("path:", path)
	if err != nil {
		fmt.Println("1111111 GetServerListByPath", err)
		return nil, err, nil
	}
	fmt.Println("zList ", zkList)
	zkConf := &LoadBalanceZkConf{
		activeList:   zkList,
		format:       format,
		zkHosts:      zkHosts,
		zkPath:       path,
		confIPWeight: conf,
	}
	zkConf.WatchConf(zkManager)
	return zkConf, nil, zkManager
}

// WatchConf 监听到服务器列表变化时，通知观察者也更新服务器列表 这里会调用updateConf 方法 去通知观察者去更新他们的服务器
func (l *LoadBalanceZkConf) WatchConf(zk *zookeeper.ZkManager) {
	//zkManager := zookeeper.NewZkManager(l.zkHosts)
	//zkManager.GetConnect()
	//fmt.Println("zkManager.WatchConf()")
	//defer zkManager.Close()
	chanList, chanErr := zk.WatchServerListByPath(l.zkPath)
	fmt.Println("WatchConf l.zkPath ", l.zkPath)
	go func() {
		for {
			select {
			case chanErr := <-chanErr:
				fmt.Println("WatchConf chanErr:", chanErr)
			case chanList := <-chanList:
				l.UpdateConf(chanList)
				fmt.Println("chanList:", chanList)
			}
		}
	}()
}

// UpdateConf 更新服务器列表 和 让每个观察者也更新服务器列表 手动更新节点列表
func (l *LoadBalanceZkConf) UpdateConf(conf []string) {
	fmt.Println("zkManager.UpdateConf()")
	l.activeList = conf
	fmt.Println("l.activeList", l.activeList)
	for _, obs := range l.observers {
		obs.Update()
	}
}

// Attach 添加观察者
func (l *LoadBalanceZkConf) Attach(o Observer) {
	l.observers = append(l.observers, o)
}

func (l *LoadBalanceZkConf) NotifyAllObservers() {
	for _, obs := range l.observers {
		obs.Update()
	}
}

// GetConf 获取服务器列表
func (l *LoadBalanceZkConf) GetConf() []string {
	var confList []string
	for _, ip := range l.activeList {
		weight, ok := l.confIPWeight[ip]
		if !ok {
			weight = "50"
		}
		confList = append(confList, fmt.Sprintf(l.format, ip)+","+weight)
	}
	return confList
}
