# `sync.Cond`
> ゴルーチンが待機したりイベントの発生を知らせるためのランデブーポイントです。

イベントが起こったことそれ自体をシグナルとして送りたい場合に用いられる (e.g. ある goroutine の処理終了).

## `sync.Cond.Signal()`
イベントが起こった場合にシグナルを送る. 起こったことしか情報がないため、条件をつけるなどして適宜補う必要がある.

### `Cond.Wait()` に関する注意
`Cond.Wait()` ループに入ると、`Cond.L.Unlock()`が呼び出され、ループから出ると`Cond.L.Lock()`が呼び出される. コードの見かけ上とは異なる動作になる.

この性質のため、`Cond.L.Lock()` -> `Cond.Wait()` ループ (アンロックとロックが行われる) -> `Cond.L.Unlock()` とする必要がある.

### `Cond.Signal` と `Cond.Wait` の例
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	con := sync.NewCond(&sync.Mutex{})  // ... (1)
	queue := make([]interface{}, 0, 10) // ... (2)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		con.L.Lock() // ... (3)
		defer con.L.Unlock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		defer con.Signal() // ... (4)
	}

	for i := 0; i < 10; i++ {
		con.L.Lock() // ... (5)
		fmt.Println("i =", i, "len =", len(queue))
		for len(queue) == 2 { // ... (6)
			fmt.Println("Waiting...")
			con.Wait() // ... (7)
			fmt.Println("End waiting")
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second) // ... (8)
		con.L.Unlock()
	}
}
```
(1) `Cond` を宣言する. `Mutex` を引数とし、 `Cond.L` でロック・アンロックを利用できる.  
(2) 操作対象の `queue` を宣言する. 常にその要素数が２となるように `Cond.Signal()`, `Cond.Wait()` を使って調整する.  
(3) `queue` を操作するために排他的アクセス権を要求. (`Cond.Wait()` ループでアンロックが呼び出される前に一度ロックする必要がある.)  
(4) `removeFromQueue()` 操作の終了時、当該 `Cond` を持つゴルーチン(今は `main()`)にシグナルを送る.  
(5) `queue` の長さを得る操作を行うため、排他的アクセス権を要求.  
(6) シグナルを待つ条件を指定して、シグナルの情報を補う.  
(7) シグナルが (4) から送られるまでゴルーチンを停止する. (`Cond.Wait()` ループに入ると、`Cond.L.Unlock()`が呼び出され、ループから出ると`Cond.L.Lock()`が呼び出される.)  
(8) ゴルーチンを追加する.

実行結果  
`queue` の長さが 2 になった場合、`removeFromQueue()` が終わるまで次の要素が追加されない.
```
i = 0 len = 0
Adding to queue
i = 1 len = 1
Adding to queue
i = 2 len = 2
Waiting...
Removed from queue
End waiting
Adding to queue
i = 3 len = 2
Waiting...
Removed from queue
End waiting
Adding to queue
i = 4 len = 2
Waiting...
Removed from queue
End waiting
Adding to queue
i = 5 len = 2
Waiting...
Removed from queue
End waiting
Adding to queue
i = 6 len = 2
Waiting...
Removed from queue
End waiting
Adding to queue
i = 7 len = 2
Waiting...
Removed from queue
End waiting
Adding to queue
i = 8 len = 2
Waiting...
Removed from queue
End waiting
Adding to queue
i = 9 len = 2
Waiting...
Removed from queue
End waiting
Adding to queue
```

## `sync.Cond.Broadcast()`
`Signal` は 内部でシグナルを待機しているゴルーチンのリストから**最も長く待っている**ゴルーチンを見つけてそのゴルーチンへシグナルを伝える.
`Broadcat` は待機しているゴルーチン**すべて**にシグナルを送る.  

### 例
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	type Button struct {
		Clicked *sync.Cond
	}
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(con *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup // ... (1)
		goroutineRunning.Add(1)
		go func() {
			defer goroutineRunning.Done()
			con.L.Lock()
			defer con.L.Unlock()
			con.Wait() // ... (2)
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(1)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	clickRegistered.Add(1)
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	clickRegistered.Add(1)
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked!")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast() // ... (3)
	clickRegistered.Wait()
}

```
(1) `subscribe` 内部で `go func` が終了する前に関数を抜けることを防ぐ `WaitGroup`  
(2) シグナルを待つ. (このとき `L.Unlock` を内部で実行する. `Wait` を抜ける時に `L.Lock` を実行する)
(3) `sync.Cond.Broadcat` で待機しているゴルーチンすべてへシグナルを送る.

実行結果
```
Mouse clicked!
Maximizing window.
Displaying annoying dialog box!
```


