package concurrenthashmap

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

/**
 *
 * @Author AiTao
 * @Date 2023/7/17 7:47
 * @Url
 **/
func TestMap_Len3(t *testing.T) {
	list := map[string]interface{}{
		"a": "aaa",
		"b": "bbb",
		"c": "ccc",
	}
	var m Map
	for k, v := range list {
		m.LoadOrStore(k, v)
	}
	fmt.Println("LoadOrStore(k,v):", m.Len(), &m)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		value, loaded := m.LoadAndDelete("a")
		fmt.Println("Delete(`a`)", m.Len(), value, loaded, &m)
		value, loaded = m.LoadAndDelete("f")
		fmt.Println("Delete(`f`)", m.Len(), value, loaded, &m)
		previous, exists := m.Swap("d", "ddd")
		fmt.Println("Store(`d`)", m.Len(), previous, exists, &m)
		value, loaded = m.LoadAndDelete("c")
		fmt.Println("Delete(`c`)", m.Len(), value, loaded, &m)
	}()
	wg.Wait()
}
func TestMap_Len2(t *testing.T) {
	m := new(Map)
	for i := 0; i < 10; i++ {
		m.LoadOrStore(i, i)
	}
	m.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})
	fmt.Println(m.Len())
}
func TestMap_Len(t *testing.T) {
	// 创建一个并发安全的 Map
	m := new(Map)

	// 设置并发读写操作的数量
	const numOps = 1000

	// 使用原子计数器来统计成功读写操作的数量
	var successOps int64

	// 使用 WaitGroup 来等待所有 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(numOps)

	// 启动多个 goroutine 进行并发读写操作
	for i := 0; i < numOps; i++ {
		go func() {
			// 生成一个随机数作为 key
			key := rand.Intn(10)

			// 进行并发的读写操作
			if rand.Float64() < 0.5 {
				m.Store(key, "value")
				atomic.AddInt64(&successOps, 1)
			} else {
				_, ok := m.Load(key)
				if ok {
					atomic.AddInt64(&successOps, 1)
				}
			}

			// 标记 goroutine 完成
			wg.Done()
		}()
	}

	// 等待所有 goroutine 完成
	wg.Wait()
	fmt.Println(m.Len())
	// 输出成功读写操作的数量
	fmt.Println("Successful operations:", successOps)
}
