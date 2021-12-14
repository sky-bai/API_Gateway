package load_balance

import (
	"errors"
)

// RoundRobinBalance 轮询负载均衡
type RoundRobinBalance struct {
	curIndex int
	rss      []string

	// 这里我需要维护下游服务器列表
	conf LoadBalanceConfInterface
}

func (r *RoundRobinBalance) Get(s string) (string, error) {
	panic("implement me")
}

func (r *RoundRobinBalance) Update() {
	panic("implement me")
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
