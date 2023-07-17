package ts_test

import (
	"fmt"
	ts "github.com/lmlat/syncmap"
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
func TestMap_Len2(t *testing.T) {
	// 实例化
	m := new(ts.Map)
	// 添加键值对
	m.Store("name", "aitao")
	m.Store("age", 100)
	// 打印键值对内容
	fmt.Println(m.String()) // {name=aitao, age=100}
	// 打印键值对数量
	fmt.Println(m.Len()) // 2
}

func TestMap_Len(t *testing.T) {
	list := map[string]interface{}{
		"a": "aaa",
		"b": "bbb",
		"c": "ccc",
	}
	var m ts.Map
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

func TestMap_Concurrency(t *testing.T) {
	m := new(ts.Map)
	// 统计成功读写操作的数量
	var count int64
	var wg sync.WaitGroup
	wg.Add(1000)
	// 启动多个 goroutine 进行并发读写操作
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			key := rand.Intn(10)
			// 进行并发的读写操作
			if rand.Float64() < 0.5 {
				m.Store(key, "value")
				atomic.AddInt64(&count, 1)
			} else {
				_, ok := m.Load(key)
				if ok {
					atomic.AddInt64(&count, 1)
				}
			}
		}()
	}
	wg.Wait()
	fmt.Println(m.Len())
	fmt.Println("成功数:", count)
}
