package load_balance

import (
	"errors"
	"strconv"
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
	conf     LoadBalanceZkConfInterface
}

func (w *WeightRoundRobinBalance) Get(s string) (string, error) {
	panic("implement me")
}

func (w *WeightRoundRobinBalance) SetConf(conf LoadBalanceZkConfInterface) {
	w.conf = conf
}

func (w *WeightRoundRobinBalance) Update() {
	panic("implement me")
}

// 每个节点的weight 是固定的 所以totalWeight 也是不变的
// 最初的时候 每个节点的effectiveWeight 和 currentWeight 分别是 weight

func (w *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("param len need 2 ")
	}

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

	return nil
}

func (w *WeightRoundRobinBalance) Next() string {
	// 我们要寻找下一个最好的节点
	var best *WeightNode
	var totalWeight int

	for i := 0; i < len(w.rss); i++ {
		node := w.rss[i]
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
		return ""
	}
	// step 5 变更best 临时权重临时权重
	best.currentWeight -= totalWeight

	return best.addr

}
