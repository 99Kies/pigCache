package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // 这个环内虚拟节点需要用到的hash算法 定义格式为：func(data []byte) uint32
	replicas int            // 每个真实节点对应的虚拟节点的个数
	keys     []int          // 所有的虚拟节点，列表，顺序排序
	hashMap  map[int]string // 虚拟节点和正式节点的映射，map[虚拟节点的hash值] => 真实节点
}

func New(replicas int, fn Hash) *Map {
	// 创建一个节点分配环
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 计算虚拟节点的 hash值.
			// strconv.Itoa(i) + key -> 说明虚拟节点的hash值大于等于hash(key)的值
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 将当前这个节点添加在环上
			m.keys = append(m.keys, hash)
			// hashMap 中增加虚拟节点和真实节点的映射关系。 map[虚拟节点hash值] => 真实节点
			m.hashMap[hash] = key
		}
	}
	// 将环上的keys进行排序
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	// 计算key的hash值
	hash := int(m.hash([]byte(key)))
	// 顺时针找到第一个匹配的虚拟节点的下标
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= (hash % m.keys[len(m.keys)-1])
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
