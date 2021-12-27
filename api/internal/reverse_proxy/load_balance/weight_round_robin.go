package load_balance

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type WeightNode struct {
	addr            string
	Weight          int // 权重值
	currentWeight   int // 当前权重值
	effectiveWeight int // 有效权重值
}

// WeightRoundRobinBalance 加权轮询负载均衡
type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*WeightNode
	rsw      int
	conf     LoadBalanceConfInterface
}

func (w *WeightRoundRobinBalance) Get(s string) (string, error) {
	fmt.Println("rss:", w.rss)
	return w.Next(), nil
}

func (w *WeightRoundRobinBalance) SetConf(conf LoadBalanceConfInterface) {
	w.conf = conf
}

func (w *WeightRoundRobinBalance) Update() {
	if conf, ok := w.conf.(*LoadBalanceCheckConfig); ok {
		//fmt.Println("WeightRoundRobinBalance get check conf:", conf.GetConf())
		w.rss = nil
		for _, ip := range conf.GetConf() {
			fmt.Println("配置负载均衡器的服务器列表WeightRoundRobinBalance get ip:", ip)
			err := w.Add(strings.Split(ip, ",")...)
			if err != nil {
				fmt.Println("WeightRoundRobinBalance get conf err:", err)
			}
		}
		fmt.Println("该负载均衡器维护的下游服务器列表", w.rss)
	}
}

// 每个节点的weight 是固定的 所以totalWeight 也是不变的
// 最初的时候 每个节点的effectiveWeight 和 currentWeight 分别是 weight

func (w *WeightRoundRobinBalance) Add(params ...string) error {
	fmt.Println("开始添加服务器")
	if len(params) != 2 {
		return errors.New("param len need 2 ")
	}

	fmt.Println("params:", params)
	fmt.Println("params:", params[1])

	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}

	// 添加节点 设置权重
	w.rss = append(w.rss, &WeightNode{
		addr:            params[0],
		Weight:          int(parInt),
		currentWeight:   int(parInt),
		effectiveWeight: int(parInt),
	})

	fmt.Println("已配置好的 WeightRoundRobinBalance get rss:", w.rss)
	return nil
}

func (w *WeightRoundRobinBalance) Next() string {
	// 我们要寻找下一个最好的节点
	var best *WeightNode
	var totalWeight int

	for i := 0; i < len(w.rss); i++ {
		node := w.rss[i]
		fmt.Println("node:", node)
		fmt.Println("node:", node.addr)
		// step1: 统计所有有效权重之和
		totalWeight += node.Weight

		// step2: 计算节点的currentWeight
		node.currentWeight += node.effectiveWeight // 计算两数之和

		// step3: 通讯异常时 -1 ,通讯成功时 +1 直到恢复到weight大小
		if node.effectiveWeight < node.Weight {
			node.effectiveWeight++
		}

		// step4: 找到最好的节点
		if best == nil || node.currentWeight > best.currentWeight {
			best = node
		}
	}

	if best == nil {
		fmt.Println("WeightRoundRobinBalance best is nil")
		return ""
	}
	// step 5 变更best 临时权重临时权重
	best.currentWeight -= totalWeight

	return best.addr

}
