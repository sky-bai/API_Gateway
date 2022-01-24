package load_balance

import (
	"errors"
	"strings"
)

// RoundRobinBalance 轮询负载均衡
type RoundRobinBalance struct {
	curIndex int
	rss      []string

	// 这里我需要维护下游服务器列表
	conf LoadBalanceConfInterface
}

func (r *RoundRobinBalance) Get(s string) (string, error) {
	return r.Next(), nil
}

func (r *RoundRobinBalance) Update() {
	if conf, ok := r.conf.(LoadBalanceConfInterface); ok {
		//fmt.Println("Update get check conf:", conf.GetConf())
		r.rss = nil
		//fmt.Println("准备更新的下游服务器列表",conf.GetConf())
		for _, ip := range conf.GetConf() {
			//fmt.Println("准备更新的单个下游服务器",strings.Split(ip, ",")[0])
			r.Add(strings.Split(ip, ",")[0])
		}
	}
}

func (r *RoundRobinBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params is empty")
	}
	r.rss = append(r.rss, params...)
	return nil
}

func (r *RoundRobinBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	if r.curIndex >= len(r.rss) {
		r.curIndex = 0
	}
	r.curIndex++
	return r.rss[r.curIndex-1]
}

func (r *RoundRobinBalance) SetConf(conf LoadBalanceConfInterface) {
	r.conf = conf
}
