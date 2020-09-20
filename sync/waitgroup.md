
# `sync.WaitGroup`
> `WaitGroup`はひとまとまりの並行処理があったとき、その結果を気にしない、もしくは他に結果を収集する手段がある場合に、それらの処理の完了を待つ手段として非常に有効です。

逆に **並行処理の結果が必要** or **結果収集手段がない** なら `select` 文を使うとよい.

## `sync.WaitGroup`の使用例
```go
var wg sync.WaitGroup

wg.Add(1) // goroutine の直前に Add() する. goroutine 内部だとスケジュールされるタイミングが定まっていないため.
go func() {
    defer wg.Done() // 必ず Done() する.
    fmt.Println("1st goroutine sleeping...")
    time.Sleep(1)
}()

wg.Add(1)
go func() {
    defer wg.Done()
    fmt.Println("2nd goroutine sleeping...")
    time.Sleep(2)
}()

wg.Wait() // main goroutine を wg が 0 になるまでブロックする.
fmt.Println("All goroutines complete.")
```

```go
package main

import (
	"fmt"
	"sync"
)

var hello = func(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	fmt.Printf("Hello from %v\n", id)
}

func main() {
	const numGreeters = 5
	var wg sync.WaitGroup
	wg.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg, i+1)
	}
	wg.Wait()
}
```