package load_balance

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// Hash  将一组数据转换成一个唯一的数据
type Hash func(data []byte) uint32
type UInt32Slice []uint32

func (u UInt32Slice) Len() int { return len(u) }
func (u UInt32Slice) Less(i, j int) bool {
	return u[i] < u[j]
}
func (u UInt32Slice) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

type ConsistentHashBalance struct {
	mux           sync.RWMutex
	hashFunc      Hash
	replicas      int               // 复制因子
	keySlice      UInt32Slice       // 以排序的节点hash 切片
	serverHashMap map[uint32]string // key 为hash值，value为节点key

	conf LoadBalanceConfInterface
}

func (c *ConsistentHashBalance) Update() {
	panic("implement me")
}

func (c *ConsistentHashBalance) IsEmpty() bool {
	if len(c.keySlice) == 0 {
		return true
	}
	return false
}

func NewConsistentHashBalance(replicas int, fn Hash) *ConsistentHashBalance {
	m := &ConsistentHashBalance{
		replicas:      replicas,
		hashFunc:      fn,
		serverHashMap: make(map[uint32]string),
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加缓冲节点，参数为节点key  切片里面保存的是服务器节点的hash值 具体是用线程安全的map 去保存
func (c *ConsistentHashBalance) Add(params ...string) error {

	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]

	fmt.Println("添加进去的 addr ", addr)

	c.mux.Lock()
	defer c.mux.Unlock() // 每次使用一个map结构的时候 都要加锁

	for i := 0; i < c.replicas; i++ {
		hash := c.hashFunc([]byte(strconv.Itoa(i) + addr))
		//fmt.Println("计算出来的 hash 值 ",hash)

		c.keySlice = append(c.keySlice, hash) // 为每一个地址添加复制因子hash出一个值 并把这个值放入一个切片里面
		c.serverHashMap[hash] = addr          // 并把这个值和地址做一个映射
	}

	fmt.Println("服务器切片的 hash 值", c.keySlice)
	// 对所有虚拟节点的哈希值进行排序，方便之后进行二分查找
	sort.Sort(c.keySlice) // 这里使用的排序是什么
	return nil
}

// Get 根据给定的对象后去获取最靠近它的那个节点
func (c *ConsistentHashBalance) Get(key string) (string, error) {
	if c.IsEmpty() == true {
		return "", errors.New("node is empty")
	}
	hash := c.hashFunc([]byte(key)) // 获取给定的对象的哈希值

	// 通过二分查找获取最优节点  第一个 服务器hash值 大于 数据hash值 就是最优服务器节点
	idx := sort.Search(len(c.keySlice), func(i int) bool { // 在服务器里面遍历 找到该key对应的hash值 找到比key对应的hash值大的服务器节点
		return c.keySlice[i] >= hash
	}) // idx 是服务器列表里面的索引值 对应一个服务器节点的hash值

	// 如果查找结果 大于 服务器节点哈希数组的最大索引，表示此时该对象哈希值位于最后一个节点之后，那么放入第一个节点
	if idx == len(c.keySlice) { // 这是二分查找找到最后一个服务器节点 就是环上第一个节点
		idx = 0
	}
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.serverHashMap[c.keySlice[idx]], nil // 根据服务器列表最优节点找到该节点的hash值 然后根据hash值 找到服务器地址

}

func (c *ConsistentHashBalance) SetConf(conf LoadBalanceConfInterface) {
	c.conf = conf
}
