package load_balance

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

type RandomBalance struct {
	curIndex int
	rss      []string

	// 这里我需要维护下游服务器列表
	conf LoadBalanceConfInterface
	// 这里定义了一个关于负载均衡配置的接口 它可以获得1.服务器配置 2.更新服务器配置 也就是说这里抽象出一组方法 不管具体实现 让其他负载均衡配置来实现这个接口
	// 每个负载均衡配置方法内部都有一个负载均衡的配置
	// 它可以获取服务器列表
}

func (r *RandomBalance) Get(s string) (string, error) {
	return r.Next(), nil
}

func (r *RandomBalance) Update() {
	if conf, ok := r.conf.(LoadBalanceConfInterface); ok {
		fmt.Println("Update get check conf:", conf.GetConf())
		r.rss = nil
		for _, ip := range conf.GetConf() {
			fmt.Println("添加进去的ip:", ip)
			r.Add(strings.Split(ip, ",")[0])
		}
		fmt.Println("该负载均衡器维护的下游服务器列表", r.rss)
	}
}

//func Obj2Json111(s interface{}) string {
//	obj, _ := jsoniter.Marshal(s)
//	return string(obj)
//}

// Add 添加服务器列表
func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params is empty")
	}

	r.rss = append(r.rss, params...)
	return nil
}

// Next 获得随机服务器IP
func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	r.curIndex = rand.Intn(len(r.rss)) //nolint:gosec
	fmt.Println("随机获取的下游服务器IP:", r.rss[r.curIndex])
	return r.rss[r.curIndex]
}

func (r *RandomBalance) SetConf(conf LoadBalanceConfInterface) {
	r.conf = conf
}

// 总结如下 :

// 这几个例子做了什么事情    这个负载均衡做的是获取服务器列表 确定下一个服务器IP

// 也就是在最初的时候 我这个结构体可以通过add方法添加服务器 next获取下一台服务器IP
// 但是在添加服务器的时候 我不想手动的去添加服务器 而是自动添加 减少服务器 它自动减少
// 我想服务器在变化的时候 就自动更新服务器列表信息

// 负载均衡做的就是一件事 确定 每一次处理请求的服务器IP

// 确定ip数组里面的值
// 这里我是已一个数组来维护服务器IP 因为数组是有序的 map是无序的
