package ts_test

import (
	"fmt"
	"github.com/lmlat/syncmap"
	"math/rand"
	"strconv"
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
		fmt.Println("Delete(`a`):", m.Len(), value, loaded, &m)
		value, loaded = m.LoadAndDelete("f")
		fmt.Println("Delete(`f`):", m.Len(), value, loaded, &m)
		previous, exists := m.Swap("d", "ddd")
		fmt.Println("Store(`d`):", m.Len(), previous, exists, &m)
		value, loaded = m.LoadAndDelete("c")
		fmt.Println("Delete(`c`):", m.Len(), value, loaded, &m)
	}()
	wg.Wait()
	fmt.Println("Current Map Size:", m.Len())
}

func TestMap_Equals(t *testing.T) {
	m := new(ts.Map)
	for i := 0; i < 10; i++ {
		m.Store(i, i)
	}
	fmt.Println("map:", m, m.Len())

	cloneMap := m.Clone()
	fmt.Println("cloneMap:", cloneMap, cloneMap.Len())

	fmt.Println("Equals:", m.Equals(cloneMap))
}

func TestMap_Concurrent_Equals(t *testing.T) {
	map1 := &ts.Map{}
	map2 := &ts.Map{}

	map1.Store("a", "aaa")
	map1.Store("b", "bbb")
	map1.Store("c", "ccc")

	wg := sync.WaitGroup{}
	concurrency := 10
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			map2.Store("a", "aaa")
			map2.Store("b", "bbb")
			map2.Store("c", "ccc")
			if map1.Equals(map2) {
				fmt.Println("goroutine - Maps are equal.")
			} else {
				fmt.Println("goroutine - Maps are not equal.")
			}
		}()
	}
	wg.Wait()
	fmt.Println("map1:", map1)
	fmt.Println("map2:", map2)
	if map1.Equals(map2) {
		fmt.Println("After concurrency test - Maps are equal.")
	} else {
		fmt.Println("After concurrency test - Maps are not equal.")
	}
}

func TestMap_Concurrency(t *testing.T) {
	m := new(ts.Map)
	var count int64
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			key := rand.Intn(10)
			m.Store(key, "value")
			_, ok := m.Load(key)
			if ok {
				atomic.AddInt64(&count, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println(m.Len()) // 键值对数量
	fmt.Println(count)   // 测试用例成功数量
}

// 测试并发访问
func TestConcurrentHashMap_Concurrent_Load(t *testing.T) {
	m := ts.Map{}
	concurrency := 100
	iterations := 1000
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				key := "key" + strconv.Itoa(j)
				value := j
				m.Store(key, value)
				result, ok := m.Load(key)
				if !ok || result != value {
					t.Errorf("Mismatch for key %s: expected %d, got %v", key, value, result)
				}
			}
		}()
	}
	wg.Wait()
}

func TestConcurrentHashMap_Concurrent_Remove(t *testing.T) {
	m := &ts.Map{}
	concurrency := 100
	iterations := 1000
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				key := "key" + strconv.Itoa(j)
				value := j
				m.Store(key, value)
				m.Delete(key)
			}
		}()
	}
	wg.Wait()
	fmt.Println("Concurrent Map Remove:", m)
}

// 测试添加多个不同的键值对
func TestMap_Swap(t *testing.T) {
	m := new(ts.Map)
	for i := 0; i < 20; i++ {
		m.Swap(i, i)
	}
	for i := 0; i < 20; i++ {
		m.Swap(i, i*10)
	}
	fmt.Println(m)
	fmt.Println(m.Len())
}

func TestMap_CompareAndSwap(t *testing.T) {
	m := new(ts.Map)
	for i := 0; i < 20; i++ {
		m.Store(i, i)
	}
	// 如果指定键的值为old, 则将其替换为new
	ok := m.CompareAndSwap(0, 0, 100) // true表示替换成功, false表示替换失败
	fmt.Println("map:", m.String(), m.Len(), ok)
}

func TestMap_LoadOrStore(t *testing.T) {
	m := new(ts.Map)
	for i := 0; i < 20; i++ {
		_, _ = m.LoadOrStore(i, i) // loaded=true表示获取了指定键的映射, loaded=false表示向Map中添加键值对, actual表示当前值
	}
	fmt.Println("map:", m.String(), m.Len())
}

func TestMap_CompareAndDelete(t *testing.T) {
	m := new(ts.Map)
	for i := 0; i < 20; i++ {
		m.LoadOrStore(i, i)
	}
	// 如果指定键的值为old则删除
	deleted := m.CompareAndDelete(2, 2)
	fmt.Println("map:", m.String(), m.Len(), deleted)
}

func TestMap_IsEmpty(t *testing.T) {
	m := new(ts.Map)
	fmt.Println(m.IsEmpty())
	for i := 0; i < 20; i++ {
		m.LoadOrStore(i, i)
	}
	fmt.Println(m.IsEmpty())
}
