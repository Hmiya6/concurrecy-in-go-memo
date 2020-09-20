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
		queue = queue[1:]
		fmt.Println("Removed from queue")
		con.L.Unlock()
		con.Signal() // ... (4)
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
