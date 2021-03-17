package component

import (
	"crypto/md5"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// 一个抽象的存储节点
type Node struct {
	Id      string
	Address string
}

// 一致性哈希
type ConsistentHash struct {
	mutex    sync.RWMutex // 读写锁
	nodes    map[int]Node // 节点
	replicas int          // 每个节点的副本数
}

func NewConsistentHash(nodes []Node, replicas int) *ConsistentHash {
	ch := &ConsistentHash{nodes: make(map[int]Node), replicas: replicas}
	for _, node := range nodes {
		ch.AddNode(node)
	}
	return ch
}

// 添加一个节点
func (ch *ConsistentHash) AddNode(node Node) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	for i := 0; i < ch.replicas; i++ {
		k := hash(node.Id + "_" + strconv.Itoa(i)) // 使用node.Id + “_” + i 计算node的hash值
		ch.nodes[k] = node                         // ch.nodes[] 中，存放了多个node的副本，位置由上述hash决定
	}
}

// 删除节点
func (ch *ConsistentHash) RemoveNode(node Node) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	for i := 0; i < ch.replicas; i++ {
		k := hash(node.Id + "_" + strconv.Itoa(i))
		delete(ch.nodes, k)
	}
}

// 传入key值，返回其所在的node节点
func (ch *ConsistentHash) GetNode(outerKey string) Node {
	key := hash(outerKey)                          // 首先求出key的hash值
	nodeKey := ch.findNearestNodeKeyClockwise(key) // 然后找到第一个比它大的nodekey，从而找到node
	return ch.nodes[nodeKey]
}

// 计算顺时针方向最近的节点
func (ch *ConsistentHash) findNearestNodeKeyClockwise(key int) int {
	ch.mutex.RLock()
	sortKeys := sortKeys(ch.nodes) // 排序后的节点key值，多个key，key和node是多对一的关系
	ch.mutex.RUnlock()
	for _, k := range sortKeys {
		if key <= k {
			return k
		}
	}
	return sortKeys[0]
}

// 返回排序后的节点key值
func sortKeys(m map[int]Node) []int {
	var sortedKeys []int
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	return sortedKeys
}

// 使用md5和crc32求hash值
func hash(key string) int {
	md5Chan := make(chan []byte, 1)
	md5Sum := md5.Sum([]byte(key))
	md5Chan <- md5Sum[:]
	return int(crc32.ChecksumIEEE(<-md5Chan))
}
