# デッドロック
> すべての平行なプロセスがお互いの処理を待ちあっている状況になっているものを指します。この状態では、プログラムは外部からの介入がない限り、決して動作する状態になりません。

デッドロックの例
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type value struct {
	mu    sync.Mutex
	value int
}

func main() {
	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock() // (1)
		defer v1.mu.Unlock() // (2)

		time.Sleep(2 * time.Second) // (3)
		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)
	go printSum(&a, &b) // (A)
	go printSum(&b, &a) // (B)
	wg.Wait()
}
```
結果
```
fatal error: all goroutines are asleep - deadlock!
```
`printSum` の説明  
(1) でクリティカルセクションへ  
(2) でクリティカルセクションを抜ける (printSumが値を返すとき)  
(3) で負荷シミュレート -> デッドロック誘発

```
main --(A)実行---(B)実行----->

(A) --- a.lock =2sec== b.lock -X
    // (A)は a をロックできるが、b をロックできない
(B) ------------ b.lock =2sec== a.lock -X
    // (B)は b をロックできるが、 a をロックできない
```

## デッドロックの条件 -- Coffman 条件
* **相互排他**  
ある並行プロセスがリソースに対して排他的な権利をどの時点においても保持している
* **条件待ち**  
ある並行プロセスはリソースの保持と追加のリソース待ちを同時に行わなければならない
* **横取り不可**  
ある並行プロセスによって保持されているリソースは、そのプロセスによってのみ解放される
* **循環待ち**  
ある並行プロセスは、他の連なっている並行プロセスを待たなければならない

これらのどれか一つでも回避できればデッドロックは発生しない