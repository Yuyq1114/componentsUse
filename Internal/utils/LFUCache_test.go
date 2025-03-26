package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestLFUCache(t *testing.T) {
	lfuCache := NewLFUCache(3) // 缓存大小为 3

	// 设置缓存
	lfuCache.Set("key1", "value1")
	lfuCache.Set("key2", "value2")
	lfuCache.Set("key3", "value3")

	// 获取缓存
	fmt.Println("获取缓存key1:", lfuCache.Get("key1")) // value1

	// 更新缓存项并更新数据库
	err := lfuCache.UpdateCacheAndDatabase("key2", "new_value2")
	if err != nil {
		fmt.Println("更新缓存失败:", err)
	}

	// 再次获取缓存
	fmt.Println("获取缓存key2:", lfuCache.Get("key2")) // new_value2

	// 添加更多缓存项，触发淘汰
	lfuCache.Set("key4", "value4")

	// 获取所有缓存项，展示淘汰机制
	fmt.Println("获取缓存key1:", lfuCache.Get("key1")) // value1
	fmt.Println("获取缓存key2:", lfuCache.Get("key2")) // new_value2
	fmt.Println("获取缓存key3:", lfuCache.Get("key3")) // null，因为它被淘汰了
	fmt.Println("获取缓存key4:", lfuCache.Get("key4")) // value4
}

// 更新缓存项，先更新数据库，再删除缓存
func (lfu *LFUCache) UpdateCacheAndDatabase(key, value string) error {
	errChan := make(chan error, 1)

	// 异步更新数据库
	go func() {
		errChan <- updateDatabase(key, value)
	}()

	// 先删除缓存
	lfu.mu.Lock()
	delete(lfu.cache, key)
	lfu.mu.Unlock()

	// 等待数据库更新结果
	if err := <-errChan; err != nil {
		return fmt.Errorf("更新数据库失败: %v", err)
	}

	// 重新设置缓存
	lfu.Set(key, value)

	return nil
}

// 模拟数据库更新函数
func updateDatabase(key, value string) error {
	// 这里模拟数据库更新，实际上可以是数据库操作代码
	time.Sleep(time.Millisecond * 100) // 模拟数据库更新时间
	fmt.Printf("数据库更新：%s = %s\n", key, value)
	return nil
}
