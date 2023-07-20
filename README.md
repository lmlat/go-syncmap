# ts.Map [![GoDoc](https://godoc.org/github.com/lmlat/go-syncmap?status.png)](https://pkg.go.dev/github.com/lmlat/go-syncmap) [![codecov](https://codecov.io/gh/lmlat/go-syncmap/branch/master/graph/badge.svg)](https://codecov.io/gh/lmlat/go-syncmap) [![Go Report Card](https://goreportcard.com/badge/github.com/lmlat/go-syncmap)](https://goreportcard.com/report/github.com/lmlat/go-syncmap)
sync.Map 是 Go 语言标准库中的一个并发安全的键值对集合, 用于在并发环境下进行读取和写入操作。
自官方在 1.9 加入了 sync.Map 之后, 就一直没有获取键值对数量的方法, 导致每次都需要调用 Range 方法来统计, 为了解决这个简单而又繁琐的问题, 基于 1.20 对 sync.Map 新增了 Len 方法, 主要用于获取 sync.Map 中存储的键值对个数。

# 实现原理
在 sync.Map 类型中增加了一个 int64 类型的属性 size, 分别在 Swap、LoadOrStore、LoadAndDelete、CompareAndDelete、dirtyLocked、missLocked方法中新增了对 size 属性的处理逻辑。

# 用法
下载:
```go
go get "github.com/lmlat/go-syncmap"
```
导入: 
```go
import "github.com/lmlat/go-syncmap"
```
注意: 导包后, 默认所有的类型都定义名为 `ts` 的包中。

示例: 
```go

// 实例化一个并发Map
m := new(ts.Map)

// 检查Map是否为空
fmt.Println(m.IsEmpty()) // {name=aitao, age=100}

// 向Map中添加键值对
m.Store("name", "aitao")
m.Store("age", 100)

// 打印键值对内容
fmt.Println(m.String()) // {name=aitao, age=100}

// 打印键值对数量
fmt.Println(m.Len()) // 2
```