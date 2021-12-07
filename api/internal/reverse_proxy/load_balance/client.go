package load_balance

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"
)

const (
	//default check setting
	DefaultCheckMethod    = 0
	DefaultCheckTimeout   = 2
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)

// LoadBalanceCheckConfig 这里还是一个管理服务器列表的对象
type LoadBalanceCheckConfig struct {
	activeList   []string          // 管理服务器列表
	confIPWeight map[string]string // 这里我要管理服务器IP和它对应的权重
	// 我还要绑定多个负载均衡的算法
	observers []Observer
	//
	format string
}

// Attach 添加观察者
func (s *LoadBalanceCheckConfig) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

// NotifyAllObservers 通知已绑定的负载均衡观察者
func (s *LoadBalanceCheckConfig) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

// GetConf 获取服务器IP和权重
func (s *LoadBalanceCheckConfig) GetConf() []string {
	return s.activeList
}

func NewLoadBalanceCheckConf(format string, conf map[string]string) (*LoadBalanceCheckConfig, error) {
	var aList []string
	for ip, _ := range conf {
		aList = append(aList, ip)
	}
	mConf := &LoadBalanceCheckConfig{
		activeList:   aList,
		confIPWeight: conf,
		format:       format,
	}
	mConf.WatchConf()
	return mConf, nil
}
func (s *LoadBalanceCheckConfig) UpdateConf(conf []string) {
	fmt.Println("UpdateConf", conf)
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

// WatchConf 配置发生变动时,通知监听者也更新
func (s *LoadBalanceCheckConfig) WatchConf() {
	fmt.Println("开始心跳检测")
	go func() {
		confIPErrNum := map[string]int{}
		for {
			changedList := []string{}
			for ip, _ := range s.confIPWeight {
				// 每一次我们去向每一台服务器发送tcp的三次握手
				conn, err := net.DialTimeout("tcp", ip, 1) // 如果是http rpc 就设置对应的方法
				// 如果成功就设置服务器故障次数为0
				if err == nil {
					conn.Close()
					if _, ok := confIPErrNum[ip]; ok {
						confIPErrNum[ip] = 0
					}

				}
				// 如果失败就设置服务器故障次数+1
				if err != nil {
					if _, ok := confIPErrNum[ip]; ok {
						confIPErrNum[ip]++
					} else {
						confIPErrNum[ip] = 1
					}
				}
				if confIPErrNum[ip] < DefaultCheckTimeout {
					changedList = append(changedList, ip)
				}
			}
			sort.Strings(changedList)
			sort.Strings(s.activeList)
			if !reflect.DeepEqual(changedList, s.activeList) {
				s.UpdateConf(changedList)
			}
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}
