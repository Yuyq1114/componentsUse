package utils

import (
	"container/heap"
	"sync"
)

// 缓存项数据结构
type CacheItem struct {
	key       string
	value     string
	frequency int
	index     int // 用于堆中的索引
}

// 优先队列，存储缓存项，按频率排序
type FrequencyHeap []*CacheItem

func (h FrequencyHeap) Len() int            { return len(h) }
func (h FrequencyHeap) Less(i, j int) bool  { return h[i].frequency < h[j].frequency }
func (h FrequencyHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *FrequencyHeap) Push(x interface{}) { *h = append(*h, x.(*CacheItem)) }
func (h *FrequencyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

type LFUCache struct {
	cache         map[string]*CacheItem
	frequencyHeap *FrequencyHeap
	mu            sync.Mutex
	maxSize       int
}

// 创建一个新的 LFUCache
func NewLFUCache(maxSize int) *LFUCache {
	h := &FrequencyHeap{}
	heap.Init(h)
	return &LFUCache{
		cache:         make(map[string]*CacheItem),
		frequencyHeap: h,
		maxSize:       maxSize,
	}
}

// 获取缓存项，若不存在则返回空字符串
func (lfu *LFUCache) Get(key string) string {
	lfu.mu.Lock()
	defer lfu.mu.Unlock()

	if item, found := lfu.cache[key]; found {
		// 更新频率
		item.frequency++
		heap.Fix(lfu.frequencyHeap, item.index)
		return item.value
	}
	return ""
}

// 设置缓存项
func (lfu *LFUCache) Set(key string, value string) {
	lfu.mu.Lock()
	defer lfu.mu.Unlock()

	// 如果缓存已满，则淘汰频率最低的项
	if len(lfu.cache) >= lfu.maxSize {
		lfu.evict()
	}

	// 添加新的缓存项
	item := &CacheItem{
		key:       key,
		value:     value,
		frequency: 1,
	}
	lfu.cache[key] = item
	heap.Push(lfu.frequencyHeap, item)
}

// 淘汰频率最低的缓存项
func (lfu *LFUCache) evict() {
	// 获取频率最小的元素
	item := heap.Pop(lfu.frequencyHeap).(*CacheItem)
	delete(lfu.cache, item.key)
}
