# `sync.Once`
 `sync.Once.Do` に渡された関数を一度だけ実行する.

## 例
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int
	increment := func() {
		count++
	}

	var once sync.Once

	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}
```

## 注意
渡された関数がなんであれ、 `once.Do` を一度だけ実行する.
```go
var once sync.Once
once.Do(doSomething1)
once.Do(doSomething2)
```
は `doSomething1` のみが実行される.