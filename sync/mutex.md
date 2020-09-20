# `sync.Mutex`

> Mutex は「相互排他」を表す "mutual exclusion" の略で、プログラムのクリティカルセクションを保護する方法の１つです.

`Mutex` は並行処理で安全な方法でこれらの共有リソースに対する排他的アクセスを提供している.  
開発者はメモリに対する慎重なアクセスを**自分で調整する**必要がある.

## 使用方法
`Mutex` を使った例.
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int
	var lock sync.Mutex // ... (1)

	increment := func() {
		lock.Lock()
		defer lock.Unlock() // ... (2)
		count++
		fmt.Printf("Incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("Decrementing: %d\n", count)
	}

	var arithmetic sync.WaitGroup
	for i := 0; i < 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done() // ... (3)
			increment()
		}()
	}

	for i := 0; i < 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	arithmetic.Wait()
	fmt.Println("Arithmetic completed.")
}
```
概要: goroutine を呼び出して `WaitGroup` で処理待ち.

(1) で `Mutex` を宣言. 変数のメモリの排他的アクセス権を得て**安全に変数を操作**する.  
(2), (3) では `defer` で unlock や done を行っているが、これは **panic 時も確実に unlock/done を行う**ため.

## `sync.RWMutex`
クリティカルセッションはプログラムのボトルネックとなる. クリティカルセッションへの出入りはコストが高いので、一般的にはクリティカルセッションで消費される時間を極力短くしようとすべき.

-> `RWMutex` は読み込み/書き込みを別にロックすることが可能. つまり、書き込みでロックがかかっていなければ、任意の数の読み込みロックが可能になる.

### 具体的な使用状況
`Mutex` ロックをかける部品のうち、読み込みのみの部分があるなら、 `RWMutex` をつかって `RWMutex.RLocker()` を呼び出して 読み込みのみのロックにするとよい.  
(大きいシステムになるほど この恩恵が大きくなる)