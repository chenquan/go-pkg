# Xbytes

# Install

```shell
go get -u github.com/chenquan/go-pkg/xbytes
```

# Use

```go
package main

import (
	"github.com/chenquan/go-pkg/xbytes"
)

func main() {
	pool := xbytes.GetNBytesPool(10)
	// Get
	b := pool.Get()
	_ = b
	// Recover 
	pool.Put(b)
}

```

# Benchmark

```shell
goos: windows
goarch: amd64
pkg: github.com/chenquan/go-pxg/xbytes
cpu: Intel(R) Core(TM) i5-8265UC CPU @ 1.60GHz
BenchmarkMakeBytes-8                2281            591449 ns/op         2264156 B/op         28 allocs/op
BenchmarkGetNBytesPool-8           10000            127145 ns/op            4811 B/op        149 allocs/op
```